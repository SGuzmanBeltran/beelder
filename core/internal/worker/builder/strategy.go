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

// GetResourceSettings returns resource settings based on plan type and server type
func GetResourceSettings(planType string, serverType string) *ResourceSettings {
	isModded := strings.Contains(strings.ToLower(serverType), "forge")

	switch planType {
	case "budget":
		if isModded {
			return &ResourceSettings{
				MemoryMin:   "1G",
				MemoryMax:   "2G",
				MemoryLimit: 2560 * OneMB, // 2.5GB
				CPULimit:    1e9,
			}
		}
		return &ResourceSettings{
			MemoryMin:   "512M",
			MemoryMax:   "1G",
			MemoryLimit: 1536 * OneMB, // 1.5GB
			CPULimit:    1e9,
		}
	case "premium":
		if isModded {
			return &ResourceSettings{
				MemoryMin:   "2G",
				MemoryMax:   "4G",
				MemoryLimit: 5 * OneMB * 1024, // 5GB
				CPULimit:    2e9,
			}
		}
		return &ResourceSettings{
			MemoryMin:   "1G",
			MemoryMax:   "2G",
			MemoryLimit: 2560 * OneMB, // 2.5GB
			CPULimit:    2e9,
		}
	default: // free
		if isModded {
			return &ResourceSettings{
				MemoryMin:   "512M",
				MemoryMax:   "1G",
				MemoryLimit: 1536 * OneMB, // 1.5GB
				CPULimit:    5e8,
			}
		}
		return &ResourceSettings{
			MemoryMin:   "512M",
			MemoryMax:   "800M",
			MemoryLimit: OneMB * 1024, // 1GB
			CPULimit:    5e8,
		}
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