package types

import "context"

// Context keys for passing metadata
type ContextKey string

const (
	// ServerIDKey is used to store server_id in context for logging traceability
	ServerIDKey ContextKey = "server_id"
)

// WithServerID adds a server ID to the context
func WithServerID(ctx context.Context, serverID string) context.Context {
	return context.WithValue(ctx, ServerIDKey, serverID)
}

// GetServerID retrieves the server ID from context
func GetServerID(ctx context.Context) (string, bool) {
	serverID, ok := ctx.Value(ServerIDKey).(string)
	return serverID, ok
}
