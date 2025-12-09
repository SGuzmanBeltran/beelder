package services

import (
	"beelder/internal/api/services/sse"
	"beelder/pkg/messaging/redpanda"
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type SSEService struct {
	consumer *redpanda.RedpandaConsumer
	hub      *sse.Hub
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewSSEService(consumerConfig *redpanda.RedpandaConsumerConfig) *SSEService {
	consumer := redpanda.NewRedpandaConsumer(consumerConfig)
	ctx, cancel := context.WithCancel(context.Background())
	return &SSEService{
		consumer: consumer,
		hub:      sse.NewHub(),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (s *SSEService) Run() {
	s.consumer.Connect()
	go s.hub.Run()
	go s.consumer.ReadMessage(s.HandleStreamMessage)
}

func (s *SSEService) Stop() error {
	s.cancel()
	s.consumer.Disconnect()
	return nil
}

func (s *SSEService) HandleStreamMessage(msg kafka.Message) (bool, error) {

	var event sse.ProgressEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return false, err
	}

	s.hub.BroadcastEvent(event)

	return true, nil
}

func (s *SSEService) GetHub() *sse.Hub {
	return s.hub
}
