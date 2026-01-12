package builder

import (
	"beelder/internal/types"
	"fmt"
	"strings"
)

// BuildStrategy defines how to build different server types
type BuildStrategy interface {
	GenerateDockerfile(config *types.CreateServerConfig) string
	GetResourceSettings() *ResourceSettings
}

// ResourceSettings holds the resource configuration for a plan
type ResourceSettings struct {
	MemoryMin   string
	MemoryMax   string
	MemoryLimit int64
	CPULimit    int64
}
const (
	OneMB = 1024 * 1024
)

// GetResourceSettings returns resource settings based on RAM plan (e.g., "1GB", "2GB", "8GB")
func GetResourceSettings(ramPlan string, serverType string) *ResourceSettings {
	// Parse RAM amount from plan (e.g., "8GB" -> 8)
	var ramGB int
	fmt.Sscanf(ramPlan, "%dGB", &ramGB)
	if ramGB == 0 {
		ramGB = 1 // Default to 1GB if parsing fails
	}

	// Calculate memory settings
	// Memory limit = chosen RAM / 4 (for demo/VPS constraints), minimum 1GB
	memoryLimitMB := max(int64(ramGB / 4 * 1024), 1024)
	memoryLimit := memoryLimitMB * OneMB

	// Max = limit - 512MB, minimum 1GB
	memoryMaxMB := max(memoryLimitMB-512, 768)

	// Min = max - 512MB, minimum 512MB
	memoryMinMB := max(memoryMaxMB-512, 512)

	return &ResourceSettings{
		MemoryMin:   fmt.Sprintf("%dM", memoryMinMB),
		MemoryMax:   fmt.Sprintf("%dM", memoryMaxMB),
		MemoryLimit: memoryLimit,
		CPULimit:    1e9, // Always 1 CPU core
	}
}

// BasicBuildStrategy for Paper/Vanilla/Fabric servers
type BasicBuildStrategy struct {
	config   *types.CreateServerConfig
	settings *ResourceSettings
}

func NewBasicBuildStrategy(config *types.CreateServerConfig) *BasicBuildStrategy {
	return &BasicBuildStrategy{
		config:   config,
		settings: GetResourceSettings(config.RamPlan, config.ServerType),
	}
}

func (s *BasicBuildStrategy) GenerateDockerfile(config *types.CreateServerConfig) string {
	return fmt.Sprintf(BasicServerTemplate, config.ServerType, s.settings.MemoryMin, s.settings.MemoryMax)
}

func (s *BasicBuildStrategy) GetResourceSettings() *ResourceSettings {
	return s.settings
}

// ForgeBuildStrategy for Forge servers
type ForgeBuildStrategy struct {
	config   *types.CreateServerConfig
	settings *ResourceSettings
}

func NewForgeBuildStrategy(config *types.CreateServerConfig) *ForgeBuildStrategy {
	return &ForgeBuildStrategy{
		config:   config,
		settings: GetResourceSettings(config.RamPlan, config.ServerType),
	}
}

func (s *ForgeBuildStrategy) GenerateDockerfile(config *types.CreateServerConfig) string {
	return fmt.Sprintf(ForgeServerTemplate, config.ServerType, s.settings.MemoryMin, s.settings.MemoryMax)
}

func (s *ForgeBuildStrategy) GetResourceSettings() *ResourceSettings {
	return s.settings
}

// StrategyFactory creates appropriate strategy
type StrategyFactory interface {
	GetStrategy(config *types.CreateServerConfig) BuildStrategy
}

type DefaultStrategyFactory struct{}

func (f *DefaultStrategyFactory) GetStrategy(config *types.CreateServerConfig) BuildStrategy {
	serverTypeLower := strings.ToLower(config.ServerType)
	switch {
	case strings.Contains(serverTypeLower, "forge"):
		return NewForgeBuildStrategy(config)
	default:
		return NewBasicBuildStrategy(config)
	}
}