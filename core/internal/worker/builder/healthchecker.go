package builder

import (
	config "beelder/internal/config/worker"
	"beelder/internal/types"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type HealthChecker struct{
	logger *slog.Logger
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		logger: slog.Default().With("component", "healthchecker"),
	}
}

// waitForServerReady checks if the Minecraft server is actually ready to accept players
// It does this by monitoring the container logs for the "Done" message that indicates server readiness
func (hc *HealthChecker) waitForServerReady(containerID string, serverData *types.CreateServerData) error {
	healthCheckerLogger := hc.logger.With(
		"server_id", serverData.ServerID,
		"server_type", serverData.ServerConfig.ServerType,
		"ram_plan", serverData.ServerConfig.RamPlan,
		"image", serverData.ImageName,
		"container_id", containerID,
	)
	healthCheckerLogger.Info("Starting health check for Minecraft server", "container_id", containerID)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.WithHost(config.WorkerEnvs.DockerHost),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to new client for healthcheck: %w", err)
	}
	defer cli.Close()

	start := time.Now()

	// Monitor container logs for the "Done" message
	for time.Since(start) < time.Duration(config.WorkerEnvs.BuilderConfig.BuildTimeout)*time.Second {
		// Get container logs
		logs, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Tail:       "200", // Get more lines to ensure we catch the message
		})
		if err != nil {
			return fmt.Errorf("failed to get container logs: %v", err)
		}

		// Read ALL available logs, not just 4096 bytes
		logContent := make([]byte, 0, 8192)
		buf := make([]byte, 4096)
		for {
			n, err := logs.Read(buf)
			if n > 0 {
				logContent = append(logContent, buf[:n]...)
			}
			if err != nil {
				break
			}
		}
		logs.Close()

		// Strip Docker log headers (first 8 bytes of each line)
		logString := hc.stripDockerHeaders(string(logContent))

		// Debug: print last few lines
		lines := strings.Split(logString, "\n")
		if len(lines) > 3 {
			healthCheckerLogger.Info("Latest logs", "lines", strings.Join(lines[len(lines)-4:], "\n"))
		}

		// Check for Minecraft server ready indicators
		if hc.isServerReady(logString) {
			healthCheckerLogger.Info("✅ Minecraft server is ready after", "duration", time.Since(start).Round(time.Second))
			return nil
		}

		// Check for common error conditions
		if hc.hasServerError(logString) {
			return fmt.Errorf("server encountered an error during startup - check logs")
		}

		// Show progress every 5 seconds
		elapsed := time.Since(start)
		if int(elapsed.Seconds())%5 == 0 {
			healthCheckerLogger.Info("⏳ Server starting...", "elapsed", elapsed.Round(time.Second))
		}

		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("timeout: Minecraft server did not become ready within %v", time.Duration(config.WorkerEnvs.BuilderConfig.BuildTimeout)*time.Second)
}

// stripDockerHeaders removes the 8-byte Docker log headers from log content
func (hc *HealthChecker) stripDockerHeaders(logContent string) string {
	lines := strings.Split(logContent, "\n")
	var cleaned []string
	for _, line := range lines {
		// Docker prepends 8 bytes: [stream_type, 0, 0, 0, size (4 bytes)]
		if len(line) > 8 {
			cleaned = append(cleaned, line[8:])
		} else if len(line) > 0 {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, "\n")
}

// isServerReady checks log content for indicators that the Minecraft server is ready
func (hc *HealthChecker) isServerReady(logContent string) bool {
	readyIndicators := []string{
		"Done",                  // Paper/Spigot: "Done (4.123s)! For help, type \"help\""
		"Server startup",          // Some servers: "Server startup complete"
		"Time elapsed:",           // Vanilla: Shows when done loading
		"For help, type \"help\"", // Common ready message
	}

	for _, indicator := range readyIndicators {
		if strings.Contains(logContent, indicator) {
			return true
		}
	}
	return false
}

// hasServerError checks for common error conditions in logs
func (hc *HealthChecker) hasServerError(logContent string) bool {
	// Only check for critical errors that actually prevent server startup
	// Ignore common warnings/errors that don't stop the server
	criticalErrorIndicators := []string{
		"Failed to bind to port",
		"OutOfMemoryError",
		"java.lang.RuntimeException",
		"Server crashed",
		"Encountered an unexpected exception",
	}

	for _, indicator := range criticalErrorIndicators {
		if strings.Contains(logContent, indicator) {
			return true
		}
	}
	return false
}