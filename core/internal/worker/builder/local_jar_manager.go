package builder

import (
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
	httpClient      HTTPClient
	logger          *slog.Logger
	assetResolver   AssetPathResolver
	versionProvider ServerVersionProvider
}

func NewLocalJarManager(assetPath string) *LocalJarManager {
	return &LocalJarManager{
		httpClient:      &http.Client{},
		logger:          slog.Default().With("component", "jar_manager"),
		assetResolver:   NewAssetResolver(assetPath),
		versionProvider: NewVersionProvider(),
	}
}

func NewLocalJarManagerWithClient(client HTTPClient, assetPath string) *LocalJarManager {
	return &LocalJarManager{
		httpClient:      client,
		logger:          slog.Default().With("component", "jar_manager"),
		assetResolver:   NewAssetResolver(assetPath),
		versionProvider: NewVersionProviderWithClient(client),
	}
}

// GetJar retrieves the JAR file for the given server configuration.
// If the JAR file is not found locally, it attempts to download it.
func (m *LocalJarManager) GetJar(ctx context.Context, serverConfig *types.CreateServerConfig) (types.JarInfo, error) {
	// Create logger with server context for traceability
	logger := m.logger.With(
		"server_type", serverConfig.ServerType,
		"version", serverConfig.ServerVersion,
	)

	// Add server_id from context if available
	if serverID, ok := types.GetServerID(ctx); ok {
		logger = logger.With("server_id", serverID)
	}

	assetsPath, err := m.assetResolver.GetAssetsPath()
	if err != nil {
		return types.JarInfo{}, fmt.Errorf("failed to resolve assets path: %w", err)
	}

	path := fmt.Sprintf("%s/executables/%s-%s.jar", assetsPath, serverConfig.ServerType, serverConfig.ServerVersion)

	_, err = os.Stat(path)
	cached := true

	if os.IsNotExist(err) {
		logger.Info("JAR not found locally, downloading...")
		
		downloadErr := m.downloadJar(ctx, path, serverConfig)

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

func (m *LocalJarManager) downloadJar(ctx context.Context, filepath string, serverConfig *types.CreateServerConfig) error {
	// Create logger with context
	logger := m.logger
	if serverID, ok := types.GetServerID(ctx); ok {
		logger = logger.With("server_id", serverID)
	}

	// Get download URL from version provider
	url, err := m.versionProvider.GetDownloadURL(serverConfig.ServerType, serverConfig.ServerVersion)
	if err != nil {
		return fmt.Errorf("failed to get download URL: %w", err)
	}

	logger.Info("Downloading JAR from URL", "url", url)

	resp, err := m.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	logger.Info("JAR downloaded successfully", "path", filepath)
	return nil
}
