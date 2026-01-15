package builder

import (
	config "beelder/internal/config/worker"
	"beelder/internal/types"
	"context"
	"encoding/json"
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
	httpClient    HTTPClient
	logger        *slog.Logger
	assetResolver AssetPathResolver
}

func NewLocalJarManager() *LocalJarManager {
	return &LocalJarManager{
		httpClient:    &http.Client{},
		logger:        slog.Default().With("component", "jar_manager"),
		assetResolver: NewAssetResolver(config.WorkerEnvs.AssetsPath),
	}
}

func NewLocalJarManagerWithClient(client HTTPClient) *LocalJarManager {
	return &LocalJarManager{
		httpClient:    client,
		logger:        slog.Default().With("component", "jar_manager"),
		assetResolver: NewAssetResolver(config.WorkerEnvs.AssetsPath),
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

	// Get download URL based on server type
	url, err := m.getDownloadURL(serverConfig)
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

func (m *LocalJarManager) getDownloadURL(serverConfig *types.CreateServerConfig) (string, error) {
	switch serverConfig.ServerType {
	case "paper":
		return m.getPaperDownloadURL(serverConfig.ServerVersion)
	case "vanilla":
		return m.getVanillaDownloadURL(serverConfig.ServerVersion)
	case "forge":
		return m.getForgeDownloadURL(serverConfig.ServerVersion)
	default:
		return "", fmt.Errorf("unsupported server type: %s", serverConfig.ServerType)
	}
}

func (m *LocalJarManager) getPaperDownloadURL(version string) (string, error) {
	// Get builds list for version
	buildsURL := fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s", version)

	resp, err := m.httpClient.Get(buildsURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch builds: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch builds: status %d", resp.StatusCode)
	}

	// Parse response to get latest build number
	var buildsData struct {
		ProjectID   string `json:"project_id"`
		ProjectName string `json:"project_name"`
		Version     string `json:"version"`
		Builds      []int  `json:"builds"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&buildsData); err != nil {
		return "", fmt.Errorf("failed to parse builds response: %w", err)
	}

	if len(buildsData.Builds) == 0 {
		return "", fmt.Errorf("no builds available for version %s", version)
	}

	// Get the latest build (last element in the array)
	latestBuild := buildsData.Builds[len(buildsData.Builds)-1]

	// Construct download URL with specific build number
	downloadURL := fmt.Sprintf(
		"https://api.papermc.io/v2/projects/paper/versions/%s/builds/%d/downloads/paper-%s-%d.jar",
		version, latestBuild, version, latestBuild,
	)

	return downloadURL, nil
}

func (m *LocalJarManager) getVanillaDownloadURL(version string) (string, error) {
	// Step 1: Get version manifest from Mojang
	manifestURL := "https://launchermeta.mojang.com/mc/game/version_manifest.json"

	resp, err := m.httpClient.Get(manifestURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch version manifest: status %d", resp.StatusCode)
	}

	// Parse manifest to find version URL
	var manifest struct {
		Latest struct {
			Release  string `json:"release"`
			Snapshot string `json:"snapshot"`
		} `json:"latest"`
		Versions []struct {
			ID          string `json:"id"`
			Type        string `json:"type"`
			URL         string `json:"url"`
			Time        string `json:"time"`
			ReleaseTime string `json:"releaseTime"`
		} `json:"versions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return "", fmt.Errorf("failed to parse version manifest: %w", err)
	}

	// Find the requested version
	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == version {
			versionURL = v.URL
			break
		}
	}

	if versionURL == "" {
		return "", fmt.Errorf("version %s not found in manifest", version)
	}

	// Step 2: Get version details to find server JAR download URL
	resp, err = m.httpClient.Get(versionURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch version details: status %d", resp.StatusCode)
	}

	// Parse version details to get server download URL
	var versionData struct {
		Downloads struct {
			Server struct {
				URL string `json:"url"`
			} `json:"server"`
		} `json:"downloads"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&versionData); err != nil {
		return "", fmt.Errorf("failed to parse version details: %w", err)
	}

	if versionData.Downloads.Server.URL == "" {
		return "", fmt.Errorf("no server download available for version %s", version)
	}

	return versionData.Downloads.Server.URL, nil
}

func (m *LocalJarManager) getForgeDownloadURL(version string) (string, error) {
	// Get promotions to find recommended/latest forge version for the minecraft version
	promotionsURL := "https://files.minecraftforge.net/net/minecraftforge/forge/promotions_slim.json"

	resp, err := m.httpClient.Get(promotionsURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch forge promotions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch forge promotions: status %d", resp.StatusCode)
	}

	var promotions struct {
		Homepage string            `json:"homepage"`
		Promos   map[string]string `json:"promos"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&promotions); err != nil {
		return "", fmt.Errorf("failed to parse forge promotions: %w", err)
	}

	// Look for recommended version first, then latest
	forgeVersion := ""
	recommendedKey := version + "-recommended"
	latestKey := version + "-latest"

	if v, ok := promotions.Promos[recommendedKey]; ok {
		forgeVersion = v
	} else if v, ok := promotions.Promos[latestKey]; ok {
		forgeVersion = v
	} else {
		return "", fmt.Errorf("no forge version found for minecraft %s (tried %s and %s)", version, recommendedKey, latestKey)
	}

	// Construct download URL for forge installer
	// Format: https://maven.minecraftforge.net/net/minecraftforge/forge/{minecraft-version}-{forge-version}/forge-{minecraft-version}-{forge-version}-installer.jar
	fullVersion := fmt.Sprintf("%s-%s", version, forgeVersion)
	downloadURL := fmt.Sprintf(
		"https://maven.minecraftforge.net/net/minecraftforge/forge/%s/forge-%s-installer.jar",
		fullVersion, fullVersion,
	)

	return downloadURL, nil
}
