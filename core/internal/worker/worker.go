package worker

import (
	config "beelder/internal/config/worker"
	"beelder/internal/worker/builder"
	"beelder/pkg/messaging/redpanda"
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

	if message.Value == nil {
		return fmt.Errorf("message value is nil")
	}
	// Process the message
	err := w.builder.BuildServer()
	if err != nil {
		fmt.Println("Error building server:", err)
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