package services

import (
	"beelder/pkg/messaging/redpanda"

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

func (s *ServerService) CreateServer() error {
	// Implement the logic to create a server
	s.producer.SendMessage(kafka.Message{
		Key:   []byte("server.create"),
		Value: []byte("New server created"),
	})
	return nil
}