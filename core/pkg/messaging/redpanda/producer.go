package redpanda

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type RedpandaConfig struct {
	Brokers []string
	Topic   string
}

type RedpandaProducer struct {
	config *RedpandaConfig
	writer *kafka.Writer
}

func NewRedpandaProducer(config *RedpandaConfig) *RedpandaProducer {
	return &RedpandaProducer{
		config: config,
	}
}

func (rp *RedpandaProducer) Connect() {
	rp.writer = &kafka.Writer{
		Addr:                   kafka.TCP(rp.config.Brokers...),
		Topic:                  rp.config.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	fmt.Println("Connected to Redpanda")
}

func (rp *RedpandaProducer) SendMessage(message kafka.Message) error {
	err := rp.writer.WriteMessages(
		context.TODO(),
		message,
	)
	if err != nil {
		return fmt.Errorf("failed to write messages: %w", err)
	}
	fmt.Println("Message sent to topic:", rp.config.Topic)
	return nil
}

func (rp *RedpandaProducer) SendJsonMessage(key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message value: %w", err)
	}

	message := kafka.Message{
		Key:   []byte(key),
		Value: jsonValue,
	}

	return rp.SendMessage(message)
}
