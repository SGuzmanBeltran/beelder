package builder

import (
	"archive/tar"
	config "beelder/internal/config/worker"
	"beelder/internal/types"
	"beelder/internal/types/interfaces"
	"beelder/pkg/messaging/redpanda"
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/docker/docker/api/types/build"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	initialPort = 25570
)

// validateServerConfig checks if the server configuration is valid before building.
// Returns an error if any required fields are missing or invalid.
func validateServerConfig(config *types.CreateServerConfig) error {
	if config.Name == "" {
		return fmt.Errorf("server name cannot be empty")
	}

	if config.ServerType == "" {
		return fmt.Errorf("server type cannot be empty")
	}

	validTypes := []string{"paper", "forge", "fabric"}
	if !contains(validTypes, config.ServerType) {
		return fmt.Errorf("invalid server type: %s (must be one of %v)", config.ServerType, validTypes)
	}

	if config.RamPlan == "" {
		return fmt.Errorf("ram plan cannot be empty")
	}

	validPlans := []string{"2GB", "4GB", "6GB", "8GB", "12GB"}
	if !contains(validPlans, config.RamPlan) {
		return fmt.Errorf("invalid ram plan: %s (must be one of %v)", config.RamPlan, validPlans)
	}

	return nil
}

// contains checks if a string slice contains a specific item.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(s, item) {
			return true
		}
	}
	return false
}

// Builder handles the creation and deployment of Minecraft servers using Docker.
// It manages the entire lifecycle from Dockerfile generation to container health checks.
type Builder struct {
	healthChecker   *HealthChecker
	producer        *redpanda.RedpandaProducer
	portCounter     atomic.Int32
	logger          *slog.Logger
	imageBuildLocks sync.Map
	jarManager      interfaces.JarManager
}

// NewBuilder initializes and returns a new Builder instance.
func NewBuilder(producer *redpanda.RedpandaProducer) *Builder {
	healthChecker := NewHealthChecker()
	builder := &Builder{
		producer:      producer,
		healthChecker: healthChecker,
		logger:        slog.Default().With("component", "builder"),
		jarManager:    NewLocalJarManager(),
	}
	builder.portCounter.Store(initialPort)
	return builder
}

// BuildServer creates and starts a new Minecraft server container.
// It generates the appropriate Dockerfile, builds the image, creates the container,
// and performs health checks before returning.
//
// Returns an error if any step fails.
func (b *Builder) BuildServer(ctx context.Context, serverData *types.CreateServerData) (error, string) {
	// Validate configuration before starting build
	if err := validateServerConfig(serverData.ServerConfig); err != nil {
		return fmt.Errorf("invalid server configuration: %w", err), "validating_configuration"
	}

	imageName := fmt.Sprintf("ms-%s-%s:latest", serverData.ServerConfig.ServerType, strings.ToLower(serverData.ServerConfig.RamPlan))
	serverData.ImageName = imageName
	builderLogger := b.logger.With(
		"action", "build_server",
		"server_id", serverData.ServerID,
		"server_type", serverData.ServerConfig.ServerType,
		"ram_plan", serverData.ServerConfig.RamPlan,
		"image", imageName,
	)
	// Strategy now handles everything including memory
	strategyFactory := &DefaultStrategyFactory{}
	buildStrategy := strategyFactory.GetStrategy(serverData.ServerConfig)
	dockerfile := buildStrategy.GenerateDockerfile(serverData.ServerConfig)

	b.producer.SendJsonMessage(
		"server.build.building",
		map[string]string{
			"message":   "Building server image",
			"status":    "building_image",
			"server_id": serverData.ServerID,
		},
	)

	err := b.buildImageFromDockerfile(ctx, dockerfile, serverData)

	if err != nil {
		return err, "building_image"
	}

	cli, err := client.NewClientWithOpts(
		client.WithHost(config.WorkerEnvs.DockerHost),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to new client: %w", err), "connecting_docker_client"
	}
	defer cli.Close()

	port := b.portCounter.Add(1) - 1
	// Create container with restart policy
	builderLogger.Info("Creating container...")
	b.producer.SendJsonMessage(
		"server.build.building",
		map[string]string{
			"message":   "Building server...",
			"status":    "building",
			"stage":     "server_creation",
			"server_id": serverData.ServerID,
		},
	)
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
						HostPort: fmt.Sprintf("%d", port),
					},
				},
			},
			RestartPolicy: container.RestartPolicy{
				Name: "unless-stopped",
			},
			Resources: container.Resources{
				Memory:   buildStrategy.GetResourceSettings().MemoryLimit,
				NanoCPUs: buildStrategy.GetResourceSettings().CPULimit,
			},
		},
		&network.NetworkingConfig{},
		nil,
		fmt.Sprintf("ms-%s-%s-%s", serverData.ServerConfig.ServerType, serverData.ServerConfig.RamPlan, serverData.ServerID),
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err), "creating_container"
	}

	serverData.ContainerID = resp.ID

	builderLogger.Info("Container created", "ID", resp.ID)

	// Start container
	b.producer.SendJsonMessage(
		"server.build.building",
		map[string]string{
			"message":   "Starting server...",
			"status":    "building",
			"stage":     "starting",
			"server_id": serverData.ServerID,
		},
	)
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err), "starting_container"
	}
	builderLogger.Info("Minecraft server started in background", "ID", resp.ID)
	builderLogger.Info("Connect to the server at localhost:25565")
	builderLogger.Info("To view logs: docker logs -f " + resp.ID)
	builderLogger.Info("To stop: docker stop " + resp.ID)

	// Health check: verify if the Minecraft server is responding on TCP port 25565
	builderLogger.Info("Checking if Minecraft server is responding...")

	b.producer.SendJsonMessage(
		"server.build.building",
		map[string]string{
			"message":   "Checking server health...",
			"status":    "building",
			"stage":     "health_checking",
			"server_id": serverData.ServerID,
		},
	)

	if err := b.healthChecker.waitForServerReady(resp.ID, serverData); err != nil {
		builderLogger.Error("health check failed, rolling back", "error", err)

		if removeErr := b.DestroyServer(ctx, serverData.ContainerID); removeErr != nil {
			return fmt.Errorf("health check failed and rollback failed: health_error=%w, remove_error=%v", err, removeErr), "health_checking"
		}

		builderLogger.Info("unhealthy container removed", "container_id", resp.ID[:12])
		return fmt.Errorf("health check failed for server %s: %w", serverData.ServerConfig.Name, err), "health_checking"
	}
	builderLogger.Info("✅ Minecraft server is ready and accepting connections!")

	return nil, "ready"
}

// DestroyServer stops and removes a Docker container by its ID.
// It forcefully removes the container if it's running.
// Returns an error if the removal fails.
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
		"ram_plan", serverData.ServerConfig.RamPlan,
		"image", serverData.ImageName,
	)

	lockInterface, _ := b.imageBuildLocks.LoadOrStore(serverData.ImageName, &sync.Mutex{})
	imageLock := lockInterface.(*sync.Mutex)

	imageLock.Lock()
	defer imageLock.Unlock()

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

	// Add server_id to context for logging traceability
	ctxWithServerID := types.WithServerID(ctx, serverData.ServerID)

	//Download the server jar
	jar, err := b.jarManager.GetJar(ctxWithServerID, serverData.ServerConfig)
	if err != nil {
		return fmt.Errorf("failed to get server jar: %w", err)
	}
	builderLogger.Info("Using server JAR", "path", jar.Path, "cached", jar.AlreadyCached)

	// Create build context
	buildContext, err := b.createBuildContext(dockerfileContent, serverData.ServerConfig.ServerType)
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
func (b *Builder) createBuildContext(dockerfileContent string, serverType string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// Add Dockerfile
	if err := addTarFile(tw, "Dockerfile", []byte(dockerfileContent)); err != nil {
		return nil, fmt.Errorf("failed to add Dockerfile to tar: %w", err)
	}

	// Use asset resolver to get assets path
	resolver := NewAssetResolver(config.WorkerEnvs.BuilderConfig.AssetsPath)
	assetsPath, err := resolver.GetAssetsPath()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve assets path: %w", err)
	}

	jarPath := filepath.Join(assetsPath, "executables", fmt.Sprintf("%s.jar", serverType))

	jarBytes, err := os.ReadFile(jarPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read server jar from %s: %w", jarPath, err)
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
