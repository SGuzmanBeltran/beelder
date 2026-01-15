package builder

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ServerVersionProvider provides information about available server versions and download URLs
type ServerVersionProvider interface {
	GetAvailableVersions(serverType string) ([]string, error)
	GetDownloadURL(serverType, version string) (string, error)
}

type DefaultVersionProvider struct {
	httpClient HTTPClient
}

func NewVersionProvider() *DefaultVersionProvider {
	return &DefaultVersionProvider{
		httpClient: &http.Client{},
	}
}

func NewVersionProviderWithClient(client HTTPClient) *DefaultVersionProvider {
	return &DefaultVersionProvider{
		httpClient: client,
	}
}

// GetAvailableVersions returns a list of available versions for the given server type
func (p *DefaultVersionProvider) GetAvailableVersions(serverType string) ([]string, error) {
	switch serverType {
	case "paper":
		return p.getPaperVersions()
	case "vanilla":
		return p.getVanillaVersions()
	case "forge":
		return p.getForgeVersions()
	default:
		return nil, fmt.Errorf("unsupported server type: %s", serverType)
	}
}

// GetDownloadURL returns the download URL for a specific server type and version
func (p *DefaultVersionProvider) GetDownloadURL(serverType, version string) (string, error) {
	switch serverType {
	case "paper":
		return p.getPaperDownloadURL(version)
	case "vanilla":
		return p.getVanillaDownloadURL(version)
	case "forge":
		return p.getForgeDownloadURL(version)
	default:
		return "", fmt.Errorf("unsupported server type: %s", serverType)
	}
}

// Paper implementation
func (p *DefaultVersionProvider) getPaperVersions() ([]string, error) {
	resp, err := p.httpClient.Get("https://api.papermc.io/v2/projects/paper")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch paper versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch paper versions: status %d", resp.StatusCode)
	}

	var data struct {
		Versions []string `json:"versions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse paper versions: %w", err)
	}

	return data.Versions, nil
}

func (p *DefaultVersionProvider) getPaperDownloadURL(version string) (string, error) {
	buildsURL := fmt.Sprintf("https://api.papermc.io/v2/projects/paper/versions/%s", version)

	resp, err := p.httpClient.Get(buildsURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch builds: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch builds: status %d", resp.StatusCode)
	}

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

	latestBuild := buildsData.Builds[len(buildsData.Builds)-1]

	downloadURL := fmt.Sprintf(
		"https://api.papermc.io/v2/projects/paper/versions/%s/builds/%d/downloads/paper-%s-%d.jar",
		version, latestBuild, version, latestBuild,
	)

	return downloadURL, nil
}

// Vanilla implementation
func (p *DefaultVersionProvider) getVanillaVersions() ([]string, error) {
	manifestURL := "https://launchermeta.mojang.com/mc/game/version_manifest.json"

	resp, err := p.httpClient.Get(manifestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch version manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch version manifest: status %d", resp.StatusCode)
	}

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
		return nil, fmt.Errorf("failed to parse version manifest: %w", err)
	}

	// Filter to only include release versions
	versions := make([]string, 0)
	for _, v := range manifest.Versions {
		if v.Type == "release" {
			versions = append(versions, v.ID)
		}
	}

	return versions, nil
}

func (p *DefaultVersionProvider) getVanillaDownloadURL(version string) (string, error) {
	manifestURL := "https://launchermeta.mojang.com/mc/game/version_manifest.json"

	resp, err := p.httpClient.Get(manifestURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch version manifest: status %d", resp.StatusCode)
	}

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

	resp, err = p.httpClient.Get(versionURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch version details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch version details: status %d", resp.StatusCode)
	}

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

// Forge implementation
func (p *DefaultVersionProvider) getForgeVersions() ([]string, error) {
	promotionsURL := "https://files.minecraftforge.net/net/minecraftforge/forge/promotions_slim.json"

	resp, err := p.httpClient.Get(promotionsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch forge promotions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch forge promotions: status %d", resp.StatusCode)
	}

	var promotions struct {
		Homepage string            `json:"homepage"`
		Promos   map[string]string `json:"promos"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&promotions); err != nil {
		return nil, fmt.Errorf("failed to parse forge promotions: %w", err)
	}

	// Extract unique minecraft versions from promo keys (e.g., "1.21.10-latest" -> "1.21.10")
	versionsMap := make(map[string]bool)
	for key := range promotions.Promos {
		// Keys are in format: "version-latest" or "version-recommended"
		var mcVersion string
		if len(key) > 7 && key[len(key)-7:] == "-latest" {
			mcVersion = key[:len(key)-7]
		} else if len(key) > 12 && key[len(key)-12:] == "-recommended" {
			mcVersion = key[:len(key)-12]
		}
		if mcVersion != "" {
			versionsMap[mcVersion] = true
		}
	}

	versions := make([]string, 0, len(versionsMap))
	for version := range versionsMap {
		versions = append(versions, version)
	}

	return versions, nil
}

func (p *DefaultVersionProvider) getForgeDownloadURL(version string) (string, error) {
	promotionsURL := "https://files.minecraftforge.net/net/minecraftforge/forge/promotions_slim.json"

	resp, err := p.httpClient.Get(promotionsURL)
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

	fullVersion := fmt.Sprintf("%s-%s", version, forgeVersion)
	downloadURL := fmt.Sprintf(
		"https://maven.minecraftforge.net/net/minecraftforge/forge/%s/forge-%s-installer.jar",
		fullVersion, fullVersion,
	)

	return downloadURL, nil
}
