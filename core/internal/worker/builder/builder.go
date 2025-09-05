package builder

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/client"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) BuildServer() error {
	_ = context.Background()

	// TODO: search the image, if not found, build it
	err := b.buildImageFromDockerfile(BasicServerTemplate, "minecraft-server:latest")

	if err != nil {
		return err
	}

	// cli, err := docker.NewClientWithOpts(
	// 	docker.WithHost("unix:///Users/card/.docker/run/docker.sock"),
	// 	docker.WithAPIVersionNegotiation(),
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// defer cli.Close()

	// _, err = cli.ImagePull(ctx, "docker.io/library/debian:bookworm-slim", image.PullOptions{})
	// if err != nil {
	// 	return err
	// }

	// // Create container
	// port, _ := nat.NewPort("tcp", "25565")
	// resp, err := cli.ContainerCreate(ctx, &container.Config{
	// 	Image: "debian:bookworm-slim",
	// 	Tty:   false,
	// 	Cmd:   []string{"sleep", "infinity"},
	// }, &container.HostConfig{
	// 	PortBindings: nat.PortMap{
	// 		port: []nat.PortBinding{
	// 			{HostIP: "0.0.0.0", HostPort: "25565"},
	// 		},
	// 	},
	// 	Resources: container.Resources{
	// 		Memory:   1 * 1024 * 1024 * 1024, // 1GB
	// 		NanoCPUs: 1 * 1e9,                // 1 CPUs
	// 	},
	// }, nil, nil, "mc-test1")

	// if err != nil {
	// 	panic(err)
	// }

	// if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Minecraft server started in container:", resp.ID[:12])

	// cmd := exec.Command("echo", "hello")

	// // Run and capture the output
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return err
	// }

	// fmt.Println("Output:", string(output))
	return nil
}

// buildImageFromDockerfile builds a Docker image from Dockerfile content (string) with the specified image name.
func (b *Builder) buildImageFromDockerfile(dockerfileContent string, imageName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.WithHost("unix:///Users/card/.docker/run/docker.sock"),
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
	buildContext, err := createBuildContext(dockerfileContent)
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
		panic(err)
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
func createBuildContext(dockerfileContent string) (io.Reader, error) {
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

	jarPath := filepath.Join(projectRoot, "assets", "executables", "paper-1.21.8-58.jar")

	jarBytes, err := os.ReadFile(jarPath)
	if err != nil {
		return nil, err
	}
	// Add to tar with the path expected by Dockerfile
	if err := addTarFile(tw, "assets/executables/paper-1.21.8-58.jar", jarBytes); err != nil {
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
