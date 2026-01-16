package types

import "net/http"

// JarInfo holds information about a JAR file.
type JarInfo struct {
	Path          string
	AlreadyCached bool
}

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}
