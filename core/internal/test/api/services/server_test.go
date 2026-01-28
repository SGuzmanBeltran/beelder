package test_services

import (
	"beelder/internal/api/services"
	"testing"
)

type MockVersionProvider struct {
	availableVersions []string
}

func (m *MockVersionProvider) GetAvailableVersions(serverType string) ([]string, error) {
	return m.availableVersions, nil
}

func (m *MockVersionProvider) GetDownloadURL(serverType, version string) (string, error) {
	return "", nil
}

func TestGetServerVersions_Filter(t *testing.T) {
	mockProvider := &MockVersionProvider{
		availableVersions: []string{
			"1.21.11",
			"1.21.11-rc2",
			"1.21.11-pre3",
			"1.20.1",
			"1.20-rc1",
		},
	}

	service := services.NewServerService(nil, mockProvider)

	versions, err := service.GetServerVersions("paper")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedCount := 2
	if len(versions) != expectedCount {
		t.Errorf("Expected %d versions, got %d", expectedCount, len(versions))
	}

	for _, v := range versions {
		if v == "1.21.11-rc2" || v == "1.21.11-pre3" || v == "1.20-rc1" {
			t.Errorf("Found filtered version: %s", v)
		}
	}
	
	// Check if stable versions exist
	foundVanilla := false
	for _, v := range versions {
		if v == "1.21.11" {
			foundVanilla = true
		}
	}
	if !foundVanilla {
		t.Error("Did not find stable version 1.21.11")
	}
}
