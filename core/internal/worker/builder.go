package worker

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/docker/go-connections/nat"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
	docker "github.com/moby/moby/client"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) BuildServer() error {
	ctx := context.Background()

	cli, err := docker.NewClientWithOpts(
		docker.WithHost("unix:///Users/card/.docker/run/docker.sock"),
		docker.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	_, err = cli.ImagePull(ctx, "docker.io/library/debian:bookworm-slim", image.PullOptions{})
	if err != nil {
		return err
	}

	// Create container
	port, _ := nat.NewPort("tcp", "25565")
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "debian:bookworm-slim",
		Tty:   false,
		Cmd:   []string{"sleep", "infinity"},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			port: []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: "25565"},
			},
		},
		Resources: container.Resources{
			Memory:   1 * 1024 * 1024 * 1024, // 1GB
			NanoCPUs: 1 * 1e9,                // 1 CPUs
		},
	}, nil, nil, "mc-test1")

	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	fmt.Println("Minecraft server started in container:", resp.ID[:12])

	cmd := exec.Command("echo", "hello")

	// Run and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	fmt.Println("Output:", string(output))
	return nil
}
