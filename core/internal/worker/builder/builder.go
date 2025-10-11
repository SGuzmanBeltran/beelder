package builder

import (
	"archive/tar"
	config "beelder/internal/config/worker"
	"beelder/internal/types"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Builder struct{
	healthChecker *HealthChecker
	portCounter int
	logger *slog.Logger
}

func NewBuilder() *Builder {
	healthChecker := NewHealthChecker()
	return &Builder{
		healthChecker: healthChecker,
		portCounter: 25565,
		logger:  slog.Default().With("component", "builder"),
	}
}

func (b *Builder) BuildServer(ctx context.Context, serverData *types.CreateServerData) error {

    imageName := fmt.Sprintf("ms-%s-%s:latest", serverData.ServerConfig.ServerType, serverData.ServerConfig.PlanType)
	serverData.ImageName = imageName
	builderLogger := b.logger.With(
		"action", "build_server",
		"server_id", serverData.ServerID,
		"server_type", serverData.ServerConfig.ServerType,
		"plan_type", serverData.ServerConfig.PlanType,
		"image", imageName,
	)
    // Strategy now handles everything including memory
    strategyFactory := &DefaultStrategyFactory{}
    buildStrategy := strategyFactory.GetStrategy(serverData.ServerConfig)
    dockerfile := buildStrategy.GenerateDockerfile(serverData.ServerConfig)

    err := b.buildImageFromDockerfile(ctx, dockerfile, serverData)

	if err != nil {
		return err
	}

	cli, err := client.NewClientWithOpts(
		client.WithHost(config.WorkerEnvs.DockerHost),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to new client: %w", err)
	}
	defer cli.Close()

	// Create container with restart policy
	builderLogger.Info("Creating container...")
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
		fmt.Sprintf("ms-%s-%s-%s", serverData.ServerConfig.ServerType, serverData.ServerConfig.PlanType, serverData.ServerID),
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	serverData.ContainerID = resp.ID

	builderLogger.Info("Container created", "ID", resp.ID)

	// Start container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	builderLogger.Info("Minecraft server started in background", "ID", resp.ID)
	builderLogger.Info("Connect to the server at localhost:25565")
	builderLogger.Info("To view logs: docker logs -f " + resp.ID)
	builderLogger.Info("To stop: docker stop " + resp.ID)

	// Health check: verify if the Minecraft server is responding on TCP port 25565
	builderLogger.Info("Checking if Minecraft server is responding...")

	if err := b.healthChecker.waitForServerReady(resp.ID, 120*time.Second, serverData); err != nil {
		builderLogger.Error("health check failed, rolling back", "error", err)

		if removeErr := b.DestroyServer(ctx, serverData.ContainerID); removeErr != nil {
			return fmt.Errorf("health check failed and rollback failed: health_error=%w, remove_error=%v", err, removeErr)
		}

		builderLogger.Info("unhealthy container removed", "container_id", resp.ID[:12])
		return fmt.Errorf("health check failed for server %s: %w", serverData.ServerConfig.Name, err)
	}
	builderLogger.Info("✅ Minecraft server is ready and accepting connections!")

	// Increment port counter for next server (if implementing multiple servers in future)
	b.portCounter++

	return nil
}

func (b *Builder) DestroyServer(ctx context.Context, containerID string) error {
	builderLogger := b.logger.With(
		"action", "destroy_server",
		"container_id", containerID,
	)
	builderLogger.Info("Destroying server...")
    cli, err := client.NewClientWithOpts(
        client.WithHost(config.WorkerEnvs.DockerHost),
    )
    if err != nil {
		builderLogger.Error("Failed to connect to Docker", "error", err)
        return fmt.Errorf("failed to connect to Docker: %w", err)
    }
    defer cli.Close()

    if err := cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
        Force: true,
    }); err != nil {
		builderLogger.Error("Failed to remove container", "error", err)
        return fmt.Errorf("failed to remove container %s: %w", containerID[:12], err)
    }

    return nil
}

// buildImageFromDockerfile builds a Docker image from Dockerfile content (string) with the specified image name.
func (b *Builder) buildImageFromDockerfile(ctx context.Context, dockerfileContent string, serverData *types.CreateServerData) error {
	builderLogger := b.logger.With(
		"server_id", serverData.ServerID,
		"server_type", serverData.ServerConfig.ServerType,
		"plan_type", serverData.ServerConfig.PlanType,
		"image", serverData.ImageName,
	)
	cli, err := client.NewClientWithOpts(
		client.WithHost(config.WorkerEnvs.DockerHost),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to new client: %w", err)
	}
	defer cli.Close()

	// Check if image already exists
	_, err = cli.ImageInspect(ctx, serverData.ImageName)
	if err == nil {
		builderLogger.Info("Image already exists, skipping build")
		return nil
	}

	builderLogger.Info("Image not found, building...")

	// Create build context
	buildContext, err := createBuildContext(dockerfileContent, serverData.ServerConfig.ServerType)
	if err != nil {
		return err
	}

	// Build the Docker image
	buildOptions := build.ImageBuildOptions{
		Tags:       []string{serverData.ImageName},
		Dockerfile: "Dockerfile",
		Remove:     true,
	}

	buildResp, err := cli.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	defer buildResp.Body.Close()

	// Optionally, print build output
	_, err = io.Copy(os.Stdout, buildResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read build output: %w", err)
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
		return nil, fmt.Errorf("failed to add Dockerfile to tar: %w", err)
	}

	// Find project root
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	jarPath := filepath.Join(projectRoot, "assets", "executables", fmt.Sprintf("%s.jar", serverType))

	jarBytes, err := os.ReadFile(jarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read server jar: %w", err)
	}
	// Add to tar with the path expected by Dockerfile
	if err := addTarFile(tw, fmt.Sprintf("assets/executables/%s.jar", serverType), jarBytes); err != nil {
		return nil, fmt.Errorf("failed to add server jar to tar: %w", err)
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
