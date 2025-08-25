package worker

import (
	config "beelder/internal/config/worker"
	"beelder/pkg/messaging/redpanda"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Worker struct {
	builder *Builder
}

func NewWorker(builder *Builder) *Worker {
	return &Worker{
		builder: builder,
	}
}

func (w *Worker) handleCreateServer(message kafka.Message) error {
	fmt.Println("Processing message:", string(message.Value))
	// Process the message
	w.builder.BuildServer()
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