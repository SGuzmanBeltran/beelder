package services

import (
	"beelder/internal/types"
	"beelder/pkg/messaging/redpanda"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type ServerService struct {
	producer *redpanda.RedpandaProducer
}

func NewServerService(brokerConfig *redpanda.RedpandaConfig) *ServerService {
	producer := redpanda.NewRedpandaProducer(brokerConfig)
	producer.Connect()
	return &ServerService{
		producer: producer,
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
