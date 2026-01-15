package builder

import (
	"fmt"
	"os"
	"path/filepath"
)

// AssetPathResolver defines how to resolve asset paths for different environments
type AssetPathResolver interface {
	GetAssetsPath() (string, error)
}

// LocalAssetResolver finds assets by walking up the directory tree (for local development)
type LocalAssetResolver struct{}

func (r *LocalAssetResolver) GetAssetsPath() (string, error) {
	projectRoot, err := r.findProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(projectRoot, "assets"), nil
}

// findProjectRoot walks up the directory tree to find the project root
// by looking for a directory that contains both "core" and "assets" subdirectories
func (r *LocalAssetResolver) findProjectRoot() (string, error) {
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
				return projectRoot, nil
			}
		}
		projectRoot = filepath.Dir(projectRoot)
	}
}

// EnvAssetResolver uses an environment variable (for Docker/production)
type EnvAssetResolver struct {
	AssetsPath string
}

func (r *EnvAssetResolver) GetAssetsPath() (string, error) {
	if r.AssetsPath == "" {
		return "", fmt.Errorf("ASSETS_PATH not configured")
	}
	return r.AssetsPath, nil
}

// NewAssetResolver creates the appropriate resolver based on configuration
func NewAssetResolver(configPath string) AssetPathResolver {
	if configPath != "" {
		return &EnvAssetResolver{AssetsPath: configPath}
	}
	return &LocalAssetResolver{}
}
