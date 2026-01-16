package test_builder

import (
	config "beelder/internal/config/worker"
	"beelder/internal/types"
	"beelder/internal/worker/builder"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Mock HTTP Client
type mockHTTPClient struct {
	response *http.Response
	err      error
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	return m.response, m.err
}

func newMockHTTPResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

// Fixtures
func newTestManager() *builder.LocalJarManager {
	return builder.NewLocalJarManager("/assets")
}

func newTestManagerWithClient(client types.HTTPClient) *builder.LocalJarManager {
	return builder.NewLocalJarManagerWithClient(client, "/assets")
}

func newTestContext() context.Context {
	return context.Background()
}

func newTestServerConfig(serverType, version string) *types.CreateServerConfig {
	return &types.CreateServerConfig{
		ServerType:    serverType,
		ServerVersion: version,
	}
}

func setupTempDir(t *testing.T) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "jar-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create executables subdirectory
	execDir := filepath.Join(tempDir, "executables")
	if err := os.MkdirAll(execDir, 0755); err != nil {
		t.Fatalf("Failed to create executables dir: %v", err)
	}

	// Set the config to use temp dir
	config.WorkerEnvs.AssetsPath = tempDir

	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	return tempDir
}

func createMockJarFile(t *testing.T, dir, serverType, version string) string {
	t.Helper()
	jarPath := filepath.Join(dir, "executables", fmt.Sprintf("%s-%s.jar", serverType, version))
	file, err := os.Create(jarPath)
	if err != nil {
		t.Fatalf("Failed to create mock jar: %v", err)
	}
	defer file.Close()

	// Write some dummy content
	file.WriteString("mock jar content")
	return jarPath
}

func TestLocalJarManager_GetJar(t *testing.T) {
	t.Run("cached jar file exists", func(t *testing.T) {
		tempDir := setupTempDir(t)
		createMockJarFile(t, tempDir, "vanilla", "1.20.1")

		manager := newTestManager()
		serverConfig := newTestServerConfig("vanilla", "1.20.1")

		jarInfo, err := manager.GetJar(newTestContext(), serverConfig)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !jarInfo.AlreadyCached {
			t.Error("expected AlreadyCached to be true")
		}
		if _, statErr := os.Stat(jarInfo.Path); statErr != nil {
			t.Errorf("JAR file does not exist at path: %s", jarInfo.Path)
		}
	})

	t.Run("download jar successfully", func(t *testing.T) {
		setupTempDir(t)

		mockClient := &mockHTTPClient{
			response: newMockHTTPResponse(http.StatusOK, "mock jar file content"),
		}
		manager := newTestManagerWithClient(mockClient)
		serverConfig := newTestServerConfig("paper", "1.19.4")

		jarInfo, err := manager.GetJar(newTestContext(), serverConfig)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if jarInfo.AlreadyCached {
			t.Error("expected AlreadyCached to be false")
		}
		if _, statErr := os.Stat(jarInfo.Path); statErr != nil {
			t.Errorf("JAR file does not exist at path: %s", jarInfo.Path)
		}
	})

	t.Run("download fails with HTTP 404", func(t *testing.T) {
		setupTempDir(t)

		mockClient := &mockHTTPClient{
			response: newMockHTTPResponse(http.StatusNotFound, "not found"),
		}
		manager := newTestManagerWithClient(mockClient)
		serverConfig := newTestServerConfig("spigot", "1.18.2")

		_, err := manager.GetJar(newTestContext(), serverConfig)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("download fails with network error", func(t *testing.T) {
		setupTempDir(t)

		mockClient := &mockHTTPClient{
			err: fmt.Errorf("network error"),
		}
		manager := newTestManagerWithClient(mockClient)
		serverConfig := newTestServerConfig("forge", "1.20.1")

		_, err := manager.GetJar(newTestContext(), serverConfig)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
