package services

import (
	"beelder/internal/helpers"
	"beelder/internal/types"
	"beelder/pkg/messaging/redpanda"
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type ServerService struct {
	producer *redpanda.RedpandaProducer
	versionProvider helpers.ServerVersionProvider
}

func NewServerService(producer *redpanda.RedpandaProducer, versionProvider helpers.ServerVersionProvider) *ServerService {
	return &ServerService{
		producer: producer,
		versionProvider: versionProvider,
	}
}

func (s *ServerService) CreateServer(serverConfig *types.CreateServerConfig) (string, error) {
	// Convert struct to JSON bytes
	serverId := uuid.New().String()
	jsonBytes, err := json.Marshal(serverConfig)
	if err != nil {
		return "", err
	}

	// Send message with JSON bytes
	go s.producer.SendMessage(kafka.Message{
		Key:   []byte("server.create"),
		Value: jsonBytes,
	})
	return serverId, nil
}

func (s *ServerService) GetRecommendedPlans(params *types.RecommendationServerParams) (types.RecommendationResponse, error) {
	// Simple recommendation logic based on players count and server type
	var plans types.RecommendationResponse

	switch params.ServerType {
	case "vanilla", "paper":
		if params.PlayerCount <= 10 {
			plans = types.RecommendationResponse{Recommendation: "2GB"}
		} else if params.PlayerCount <= 30 {
			plans = types.RecommendationResponse{Recommendation: "4GB"}
		} else if params.PlayerCount <= 50 {
			plans = types.RecommendationResponse{Recommendation: "6GB"}
		} else if params.PlayerCount <= 100 {
			plans = types.RecommendationResponse{Recommendation: "8GB"}
		}
	case "forge":
		if params.PlayerCount <= 10 {
			plans = types.RecommendationResponse{Recommendation: "4GB"}
		} else if params.PlayerCount <= 30 {
			plans = types.RecommendationResponse{Recommendation: "6GB"}
		} else if params.PlayerCount <= 50 {
			plans = types.RecommendationResponse{Recommendation: "8GB"}
		} else if params.PlayerCount <= 100 {
			plans = types.RecommendationResponse{Recommendation: "12GB"}
		}
	}

	return plans, nil
}

func (s *ServerService) GetServerVersions(serverType string) ([]string, error) {
	supportedServerTypes := []string{"vanilla", "paper", "forge"}

	if !slices.Contains(supportedServerTypes, serverType) {
		return nil, fmt.Errorf("unsupported server type: %s", serverType)
	}

	versions, err := s.versionProvider.GetAvailableVersions(serverType)

	if err != nil {
		return nil, err
	}

	// Filter out experimental versions (containing "-")
	// This cleans up lists like Paper which include RC/Pre-release versions
	var cleanVersions []string
	for _, v := range versions {
		if !strings.Contains(v, "-") {
			cleanVersions = append(cleanVersions, v)
		}
	}
	versions = cleanVersions

	// Sort versions descending (newest first)
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i], versions[j])
	})

	return versions, nil
}

// Helper function to compare versions
func compareVersions(v1, v2 string) bool {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := range maxLen {
		var n1, n2 int
		if i < len(parts1) {
			n1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			n2, _ = strconv.Atoi(parts2[i])
		}

		if n1 != n2 {
			return n1 > n2
		}
	}
	return false
}
