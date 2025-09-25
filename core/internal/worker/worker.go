package worker

import (
	config "beelder/internal/config/worker"
	"beelder/internal/types"
	"beelder/internal/worker/builder"
	"beelder/pkg/messaging/redpanda"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Worker struct {
	builder *builder.Builder
}

func NewWorker(builder *builder.Builder) *Worker {
	return &Worker{
		builder: builder,
	}
}

func (w *Worker) handleCreateServer(message kafka.Message) error {
	fmt.Println("Processing message:", string(message.Value))

	// Create a pointer to the config struct
    serverConfig := &types.CreateServerConfig{}

    // Unmarshal JSON into the struct
    if err := json.Unmarshal(message.Value, serverConfig); err != nil {
        return fmt.Errorf("failed to unmarshal server config: %w", err)
    }

	fmt.Printf("Processed config: %+v\n", serverConfig)

	if message.Value == nil {
		return fmt.Errorf("message value is nil")
	}

    // Pass the config to BuildServer
    if err := w.builder.BuildServer(serverConfig); err != nil {
        return fmt.Errorf("failed to build server: %w", err)
    }

	return nil
}

func (w *Worker) Start() error {
	// Implement the logic to start the worker
	redpandaConsumer := redpanda.NewRedpandaConsumer(&redpanda.RedpandaConsumerConfig{
		Brokers: []string{config.WorkerEnvs.Broker},
		Topic:   config.WorkerEnvs.ConsumerTopic,
		GroupID: config.WorkerEnvs.GroupID,
	})
	redpandaConsumer.Connect()
	defer redpandaConsumer.Disconnect()

	fmt.Println("Worker started")
	redpandaConsumer.ReadMessage(w.handleCreateServer)
	fmt.Println("Closing worker")

	return nil
}