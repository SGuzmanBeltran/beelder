package main

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
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

var dockerfileTemplate = `FROM alpine:latest

# Install necessary packages (openjdk for Minecraft, bash, curl, etc.)
RUN apk add --no-cache openjdk21-jre bash curl

# Set working directory
WORKDIR /server

# Copy the Minecraft server jar
COPY assets/executables/paper-1.21.8-56.jar /server/server.jar

# Expose default Minecraft port
EXPOSE 25565

# Accept EULA by default
RUN echo "eula=true" > eula.txt

# Default command to start the Minecraft server
CMD ["java", "-Xmx1024M", "-Xms1024M", "-jar", "server.jar", "nogui"]
`

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	imageName := "minecraft-server:latest"

	// Create build context
	buildContext, err := createBuildContext()
	if err != nil {
		panic(err)
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

	// Stream build output
	io.Copy(os.Stdout, buildResp.Body)

	// Create container with restart policy
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
						HostPort: "25565",
					},
				},
			},
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
		},
		&network.NetworkingConfig{},
		nil,
		"minecraft-server",
	)
	if err != nil {
		panic(err)
	}

	// Start container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}
	fmt.Printf("Minecraft server started in background with ID: %s\n", resp.ID)
	fmt.Println("Connect to the server at localhost:25565")
	fmt.Printf("To view logs: docker logs -f %s\n", resp.ID)
	fmt.Printf("To stop: docker stop %s\n", resp.ID)

	// Optional: Show initial logs for a few seconds then exit
	fmt.Println("Showing initial startup logs...")
	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false, // Don't follow, just get initial logs
		Tail:       "50",
	})
	if err == nil {
		defer out.Close()
		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}
}

// createBuildContext creates a tar archive in memory with the Dockerfile and the server jar
func createBuildContext() (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()
	// Add Dockerfile
	if err := addTarFile(tw, "Dockerfile", []byte(dockerfileTemplate)); err != nil {
		return nil, err
	}
	// Add server jar
	jarPath := "assets/executables/paper-1.21.8-56.jar"
	jarBytes, err := os.ReadFile(jarPath)
	if err != nil {
		return nil, err
	}
	if err := addTarFile(tw, jarPath, jarBytes); err != nil {
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