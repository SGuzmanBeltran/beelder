package redpanda

import (
	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

var logger = slog.Default()

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
		CommitInterval: 0,
	})

	logger.Info("Connected to Redpanda")
}

func (rc *RedpandaConsumer) Disconnect() {
	if err := rc.reader.Close(); err != nil {
		logger.Error("Error closing reader", "error", err)
	} else {
		logger.Info("Reader closed")
	}
}

func (rc *RedpandaConsumer) ReadMessage(callback func(kafka.Message) (bool, error)) error {
	ctx := context.Background()
	for {
		m, err := rc.reader.FetchMessage(ctx)
		if err != nil {
			logger.Error("Error reading message", "error", err)
			continue
		}

		logger.Info("Message received", "topic", m.Topic, "partition", m.Partition, "offset", m.Offset, "key", string(m.Key), "value", string(m.Value))

		// Process message in goroutine for concurrency
		go func(msg kafka.Message) {
			commit, err := callback(msg)
			if err != nil {
				logger.Error("processing failed", "error", err)
			}

			if !commit {
				return // Don't commit this message
			}

			if err := rc.reader.CommitMessages(ctx, msg); err != nil {
				logger.Error("failed to commit messages", "error", err)
			}
		}(m)
	}
}
