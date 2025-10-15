package worker

import (
	config "beelder/internal/config/worker"
	"beelder/internal/types"
	"beelder/internal/worker/builder"
	"beelder/pkg/messaging/redpanda"
	"context"
	"encoding/json"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

// Worker represents a worker that processes messages from a message broker,
// builds servers, and manages concurrency limits.
type Worker struct {
	producer *redpanda.RedpandaProducer
	builder *builder.Builder
	logger  *slog.Logger
	currentServerBuilds 	atomic.Int32
	currentLiveServers   	atomic.Int32
}

// NewWorker creates and returns a new Worker instance with initialized components.
func NewWorker() *Worker {
	producer := redpanda.NewRedpandaProducer(&redpanda.RedpandaConfig{
			Brokers: []string{config.WorkerEnvs.Broker},
			Topic:   config.WorkerEnvs.ProducerTopic,
		})
	producer.Connect()
	return &Worker{
		builder: builder.NewBuilder(),
		producer: producer,
		logger:  slog.Default().With("component", "worker"),
	}
}

// handleCreateServer processes a "server.create" message.
// It checks for concurrency limits, builds the server, and sends success or failure messages.
//
// Returns a boolean indicating whether the message should be commited or not and an error if any occurred.
func (w *Worker) handleCreateServer(message kafka.Message) (bool, error) {
	if w.currentServerBuilds.Load() >= config.WorkerEnvs.BuilderConfig.MaxConcurrentBuilds {
		w.logger.Warn("Max concurrent server builds reached, skipping message")
		time.Sleep(5 * time.Second)  // Wait before retrying
		return false, nil
	}

	if w.currentLiveServers.Load() >= config.WorkerEnvs.BuilderConfig.MaxAliveServers {
		w.logger.Warn("Max alive servers reached, skipping message")
		time.Sleep(5 * time.Second)  // Wait before retrying
		return false, nil
	}

	w.currentServerBuilds.Add(1)
	defer w.currentServerBuilds.Add(-1)

	ctx := context.Background()
	serverId := uuid.New().String()
	createLogger := w.logger.With(
		"server_id", serverId,
	)
	createLogger.Info("Received create server message", "Value", string(message.Value))

	serverConfig := &types.CreateServerConfig{}
	if err := json.Unmarshal(message.Value, serverConfig); err != nil {
		createLogger.Error("Failed to unmarshal server config", "error", err)
		w.producer.SendMessage(kafka.Message{
			Key:   []byte("server.create.failed"),
			Value: []byte(err.Error()),
		})
		return true, err
	}

	createServerData := &types.CreateServerData{
		ServerID:     serverId,
		ServerConfig: serverConfig,
	}

	createLogger = createLogger.With(
		"server_type", serverConfig.ServerType,
		"plan_type", serverConfig.PlanType,
	)

	createLogger.Info("building server")
	if err := w.builder.BuildServer(ctx, createServerData); err != nil {
		createLogger.Error("server build failed", "error", err)
		w.producer.SendMessage(kafka.Message{
			Key:   []byte("server.create.failed"),
			Value: []byte(err.Error()),
		})
		return true, err
	}

	w.currentLiveServers.Add(1)
	createLogger.Info("server created successfully")
	w.producer.SendMessage(kafka.Message{
		Key:   []byte("server.create.success"),
		Value: []byte("Server created successfully"),
	})
	return true, nil
}

// handleMessage processes incoming Kafka messages and routes them to the appropriate handler based on the message key.
// It returns a boolean indicating whether the message should be committed or not and an error if any occurred.
func (w *Worker) handleMessage(message kafka.Message) (bool, error) {
	// Implement the logic to handle incoming messages
	w.logger.Info("Received message", "Value", string(message.Value))

	msgType := message.Key
	switch string(msgType) {
	case "server.create":
		return w.handleCreateServer(message)
	default:
		w.logger.Warn("Unknown message type", "type", string(msgType))
	}

	return true, nil
}

// Start initializes the worker, connects to the message broker, and begins processing messages.
func (w *Worker) Start() error {
	// Implement the logic to start the worker
	redpandaConsumer := redpanda.NewRedpandaConsumer(&redpanda.RedpandaConsumerConfig{
		Brokers: []string{config.WorkerEnvs.Broker},
		Topic:   config.WorkerEnvs.ConsumerTopic,
		GroupID: config.WorkerEnvs.GroupID,
	})
	redpandaConsumer.Connect()
	defer redpandaConsumer.Disconnect()

	w.logger.Info("Worker started and listening for messages")
	redpandaConsumer.ReadMessage(w.handleMessage)
	w.logger.Info("Closing worker")

	return nil
}
