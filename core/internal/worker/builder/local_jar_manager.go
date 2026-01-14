package builder

import (
	config "beelder/internal/config/worker"
	"beelder/internal/types"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type LocalJarManager struct {
	httpClient HTTPClient
	logger     *slog.Logger
}

func NewLocalJarManager() *LocalJarManager {
	return &LocalJarManager{
		httpClient: &http.Client{},
		logger:     slog.Default().With("component", "jar_manager"),
	}
}

func NewLocalJarManagerWithClient(client HTTPClient) *LocalJarManager {
	return &LocalJarManager{
		httpClient: client,
		logger:     slog.Default().With("component", "jar_manager"),
	}
}

// GetJar retrieves the JAR file for the given server configuration.
// If the JAR file is not found locally, it attempts to download it.
func (m *LocalJarManager) GetJar(ctx context.Context, serverConfig *types.CreateServerConfig) (types.JarInfo, error) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return types.JarInfo{}, fmt.Errorf("failed to find assets path: %w (set ASSETS_PATH env var)", err)
	}
	path := fmt.Sprintf("%s/%s/executables/%s-%s.jar", projectRoot, config.WorkerEnvs.AssetsPath, serverConfig.ServerType, serverConfig.ServerVersion)

	_, err = os.Stat(path)
	cached := true

	if os.IsNotExist(err) {
		m.logger.Info("JAR not found locally, downloading...", "server_type", serverConfig.ServerType, "version", serverConfig.ServerVersion)
		url := fmt.Sprintf("https://mcutils.com/api/server-jars/%s/%s/download", serverConfig.ServerType, serverConfig.ServerVersion)
		downloadErr := m.downloadJar(path, url)

		if downloadErr != nil {
			return types.JarInfo{}, fmt.Errorf("JAR file not found at %s and failed to download: %w", path, downloadErr)
		}

		cached = false
		err = nil // Clear the error after successful download
	}

	if err != nil {
		return types.JarInfo{}, fmt.Errorf("failed to access JAR file at %s: %w", path, err)
	}

	return types.JarInfo{
		Path:          path,
		AlreadyCached: cached,
	}, nil
}

func (m *LocalJarManager) downloadJar(filepath string, url string) error {
	resp, err := m.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file using io.Copy
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	m.logger.Info("JAR downloaded successfully", "path", filepath)
	return nil
}
