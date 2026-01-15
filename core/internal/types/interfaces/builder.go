package interfaces

import (
	"beelder/internal/types"
	"context"
)

// JarManager defines the contract for managing Minecraft server JAR files.
type JarManager interface {
	// GetJar retrieves a Minecraft server JAR file for the specified configuration.
	// Returns an error if the download fails or the serveri type/version is invalid.
	GetJar(ctx context.Context, serverConfig *types.CreateServerConfig) (types.JarInfo, error)
}
