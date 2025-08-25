package redpanda

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type RedpandaConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type RedpandaConsumer struct {
	config *RedpandaConsumerConfig
	reader *kafka.Reader
}

func NewRedpandaConsumer(config *RedpandaConsumerConfig) *RedpandaConsumer {
	return &RedpandaConsumer{
		config: config,
	}
}

func (rc *RedpandaConsumer) Connect() {
	rc.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        rc.config.Brokers,
		Topic:          rc.config.Topic,
		GroupID:        rc.config.GroupID,
		CommitInterval: time.Second,
	})

	fmt.Println("Connected to Redpanda")
}

func (rc *RedpandaConsumer) Disconnect() {
	if err := rc.reader.Close(); err != nil {
		log.Println("Error closing reader:", err)
	} else {
		fmt.Println("Reader closed")
	}
}

func (rc *RedpandaConsumer) ReadMessage(callback func(kafka.Message) error) {
	ctx := context.Background()
	for {
		m, err := rc.reader.FetchMessage(ctx)
		if err != nil {
			fmt.Println("Error reading message:", err)
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		if err := callback(m); err != nil {
			log.Println("processing failed, not committing offset:", err)
			// optionally push to DLQ or backoff, then continue
			continue
		}

		if err := rc.reader.CommitMessages(ctx, m); err != nil {
			log.Fatal("failed to commit messages:", err)
		}
	}
}
