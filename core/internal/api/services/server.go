package services

import (
	"beelder/internal/types"
	"beelder/pkg/messaging/redpanda"
	"encoding/json"

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

func (s *ServerService) CreateServer(serverConfig *types.CreateServerConfig) error {
	// Convert struct to JSON bytes
	jsonBytes, err := json.Marshal(serverConfig)
    if err != nil {
        return err
    }

	// Send message with JSON bytes
	go s.producer.SendMessage(kafka.Message{
		Key:   []byte("server.create"),
		Value: jsonBytes,
	})
	return nil
}