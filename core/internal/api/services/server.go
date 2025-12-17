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