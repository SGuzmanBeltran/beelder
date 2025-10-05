package builder

import (
	"archive/tar"
	config "beelder/internal/config/worker"
	"beelder/internal/types"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type Builder struct{
	healthChecker *HealthChecker
	portCounter int
}

func NewBuilder() *Builder {
	healthChecker := NewHealthChecker()
	return &Builder{
		healthChecker: healthChecker,
		portCounter: 25565,
	}
}

func (b *Builder) BuildServer(serverConfig *types.CreateServerConfig) error {
    ctx := context.Background()

    imageName := fmt.Sprintf("ms-%s-%s:latest", serverConfig.ServerType, serverConfig.PlanType)

    // Strategy now handles everything including memory
    strategyFactory := &DefaultStrategyFactory{}
    buildStrategy := strategyFactory.GetStrategy(serverConfig)
    dockerfile := buildStrategy.GenerateDockerfile(serverConfig)

    err := b.buildImageFromDockerfile(dockerfile, imageName, serverConfig)

	if err != nil {
		return err
	}

	cli, err := client.NewClientWithOpts(
		client.WithHost(config.WorkerEnvs.DockerHost),
	)
	if err != nil {
		return err
	}
	defer cli.Close()

	// Create container with restart policy
	fmt.Println("Creating container...")
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
			ExposedPorts: nat.PortSet{
				"25565/tcp": {},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"25565/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: fmt.Sprintf("%d", b.portCounter),
					},
				},
			},
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		},
		&network.NetworkingConfig{},
		nil,
		fmt.Sprintf("ms-%s-%s-%s", serverConfig.ServerType, serverConfig.PlanType, uuid.New().String()),
	)
	if err != nil {
		return err
	}

	fmt.Println("Container created with ID:", resp.ID)
	// Start container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}
	fmt.Printf("Minecraft server started in background with ID: %s\n", resp.ID)
	fmt.Println("Connect to the server at localhost:25565")
	fmt.Printf("To view logs: docker logs -f %s\n", resp.ID)
	fmt.Printf("To stop: docker stop %s\n", resp.ID)

	// Health check: verify if the Minecraft server is responding on TCP port 25565
	fmt.Println("Checking if Minecraft server is responding...")

	if err := b.healthChecker.waitForServerReady(resp.ID, 120*time.Second); err != nil {
		fmt.Printf("Warning: Server health check failed: %v\n", err)
		fmt.Println("The server might still be starting up. Check logs with: docker logs -f", resp.ID)
	} else {
		fmt.Println("✅ Minecraft server is ready and accepting connections!")
	}

	// Increment port counter for next server (if implementing multiple servers in future)
	b.portCounter++

	return nil
}

// buildImageFromDockerfile builds a Docker image from Dockerfile content (string) with the specified image name.
func (b *Builder) buildImageFromDockerfile(dockerfileContent string, imageName string, serverConfig *types.CreateServerConfig) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.WithHost(config.WorkerEnvs.DockerHost),
	)
	if err != nil {
		return err
	}
	defer cli.Close()

	// Check if image already exists
	_, err = cli.ImageInspect(ctx, imageName)
	if err == nil {
		fmt.Printf("Image %s already exists, skipping build\n", imageName)
		return nil
	}

	fmt.Printf("Image %s not found, building...\n", imageName)

	// Create build context
	buildContext, err := createBuildContext(dockerfileContent, serverConfig.ServerType)
	if err != nil {
		return err
	}

	// Build the Docker image
	buildOptions := build.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: "Dockerfile",
		Remove:     true,
	}

	buildResp, err := cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return err
	}
	defer buildResp.Body.Close()

	// Optionally, print build output
	_, err = io.Copy(os.Stdout, buildResp.Body)
	if err != nil {
		return err
	}

	return nil
}

// createBuildContext creates a tar archive in memory with the Dockerfile and the server jar
func createBuildContext(dockerfileContent string, serverType string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	// Add Dockerfile
	if err := addTarFile(tw, "Dockerfile", []byte(dockerfileContent)); err != nil {
		return nil, err
	}

	// Find project root
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	jarPath := filepath.Join(projectRoot, "assets", "executables", fmt.Sprintf("%s.jar", serverType))

	jarBytes, err := os.ReadFile(jarPath)
	if err != nil {
		return nil, err
	}
	// Add to tar with the path expected by Dockerfile
	if err := addTarFile(tw, fmt.Sprintf("assets/executables/%s.jar", serverType), jarBytes); err != nil {
		return nil, err
	}
	return buf, nil
}

// addTarFile adds a file to a tar writer
func addTarFile(tw *tar.Writer, name string, data []byte) error {
	hdr := &tar.Header{
		Name:    filepath.ToSlash(name),
		Mode:    0644,
		Size:    int64(len(data)),
		ModTime: time.Now(),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err := tw.Write(data)
	return err
}

// findProjectRoot walks up the directory tree to find the project root
// by looking for a directory that contains both "core" and "assets" subdirectories
func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	projectRoot := wd
	for {
		if filepath.Base(projectRoot) == "/" {
			return "", fmt.Errorf("could not find project root")
		}
		coreDir := filepath.Join(projectRoot, "core")
		assetsDir := filepath.Join(projectRoot, "assets")
		if _, err := os.Stat(coreDir); err == nil {
			if _, err := os.Stat(assetsDir); err == nil {
				break // Found project root
			}
		}
		projectRoot = filepath.Dir(projectRoot)
	}

	return projectRoot, nil
}
