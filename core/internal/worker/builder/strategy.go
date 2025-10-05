package builder

import (
	"beelder/internal/types"
	"fmt"
	"strings"
)

// BuildStrategy defines how to build different server types
type BuildStrategy interface {
	GenerateDockerfile(config *types.CreateServerConfig) string
	GetMemorySettings(planType string) *types.MemorySettings
}

// BasicBuildStrategy for Paper/Vanilla servers
type BasicBuildStrategy struct{}

func (s *BasicBuildStrategy) GenerateDockerfile(config *types.CreateServerConfig) string {
	memory := s.GetMemorySettings(config.PlanType)
	return fmt.Sprintf(BasicServerTemplate, config.ServerType, memory.Min, memory.Max)
}

func (s *BasicBuildStrategy) GetMemorySettings(planType string) *types.MemorySettings {
	switch planType {
	case "budget":
		return &types.MemorySettings{Min: "512M", Max: "1G"}
	case "premium":
		return &types.MemorySettings{Min: "1G", Max: "2G"}
	default: // free
		return &types.MemorySettings{Min: "512M", Max: "512M"}
	}
}

// ForgeBuildStrategy for Forge servers
type ForgeBuildStrategy struct{}

func (s *ForgeBuildStrategy) GenerateDockerfile(config *types.CreateServerConfig) string {
	memory := s.GetMemorySettings(config.PlanType)
	return fmt.Sprintf(ForgeServerTemplate, config.ServerType, memory.Min, memory.Max)
}

func (s *ForgeBuildStrategy) GetMemorySettings(planType string) *types.MemorySettings {
	// Forge might need different memory allocations
	switch planType {
	case "budget":
		return &types.MemorySettings{Min: "1G", Max: "2G"}
	case "premium":
		return &types.MemorySettings{Min: "2G", Max: "4G"}
	default:
		return &types.MemorySettings{Min: "512M", Max: "1G"}
	}
}

// FabricBuildStrategy for Fabric servers
type FabricBuildStrategy struct{}

func (s *FabricBuildStrategy) GenerateDockerfile(config *types.CreateServerConfig) string {
	memory := s.GetMemorySettings(config.PlanType)
	return fmt.Sprintf(BasicServerTemplate, config.ServerType, memory.Min, memory.Max)
}

func (s *FabricBuildStrategy) GetMemorySettings(planType string) *types.MemorySettings {
	switch planType {
	case "budget":
		return &types.MemorySettings{Min: "512M", Max: "1G"}
	case "premium":
		return &types.MemorySettings{Min: "1G", Max: "2G"}
	default:
		return &types.MemorySettings{Min: "512M", Max: "512M"}
	}
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
		return &ForgeBuildStrategy{}
	case strings.Contains(serverTypeLower, "fabric"):
		return &FabricBuildStrategy{}
	default:
		return &BasicBuildStrategy{}
	}
}