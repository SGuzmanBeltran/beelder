package builder

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type HealthChecker struct{}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

// waitForServerReady checks if the Minecraft server is actually ready to accept players
// It does this by monitoring the container logs for the "Done" message that indicates server readiness
func (hc *HealthChecker) waitForServerReady(containerID string, timeout time.Duration) error {
	fmt.Printf("Waiting for Minecraft server to be ready (timeout: %v)...\n", timeout)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.WithHost("unix:///Users/card/.docker/run/docker.sock"),
	)
	if err != nil {
		return err
	}
	defer cli.Close()

	start := time.Now()

	// Monitor container logs for the "Done" message
	for time.Since(start) < timeout {
		// Get container logs from the last few seconds
		logs, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Tail:       "50", // Get last 50 lines
		})
		if err != nil {
			return fmt.Errorf("failed to get container logs: %v", err)
		}

		// Read the logs
		logContent := make([]byte, 4096)
		n, _ := logs.Read(logContent)
		logs.Close()

		logString := string(logContent[:n])

		// Check for Minecraft server ready indicators
		if hc.isServerReady(logString) {
			fmt.Printf("✅ Minecraft server is ready after %v\n", time.Since(start).Round(time.Second))
			return nil
		}

		// Check for common error conditions
		if hc.hasServerError(logString) {
			return fmt.Errorf("server encountered an error during startup - check logs")
		}

		// Show progress every 5 seconds
		elapsed := time.Since(start)
		if int(elapsed.Seconds())%5 == 0 {
			fmt.Printf("⏳ Server starting... (%v elapsed)\n", elapsed.Round(time.Second))
		}

		time.Sleep(2 * time.Second)
	}

	return fmt.Errorf("timeout: Minecraft server did not become ready within %v", timeout)
}

// isServerReady checks log content for indicators that the Minecraft server is ready
func (hc *HealthChecker) isServerReady(logContent string) bool {
	readyIndicators := []string{
		"Done (",                  // Paper/Spigot: "Done (4.123s)! For help, type \"help\""
		"Server startup",          // Some servers: "Server startup complete"
		"Time elapsed:",           // Vanilla: Shows when done loading
		"For help, type \"help\"", // Common ready message
	}

	for _, indicator := range readyIndicators {
		if contains(logContent, indicator) {
			return true
		}
	}
	return false
}

// hasServerError checks for common error conditions in logs
func (hc *HealthChecker) hasServerError(logContent string) bool {
	errorIndicators := []string{
		"Exception",
		"Error",
		"Failed to bind to port",
		"Could not load",
		"OutOfMemoryError",
	}

	for _, indicator := range errorIndicators {
		if contains(logContent, indicator) {
			return true
		}
	}
	return false
}


// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				containsAt(s, substr, 1)))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}