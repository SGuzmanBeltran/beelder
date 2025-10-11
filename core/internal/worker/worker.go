package worker

import (
	config "beelder/internal/config/worker"
	"beelder/internal/types"
	"beelder/internal/worker/builder"
	"beelder/pkg/messaging/redpanda"
	"context"
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type Worker struct {
	builder *builder.Builder
	logger  *slog.Logger
}

func NewWorker(builder *builder.Builder) *Worker {
	return &Worker{
		builder: builder,
		logger:  slog.Default().With("component", "worker"),
	}
}

func (w *Worker) handleCreateServer(message kafka.Message) error {
	ctx := context.Background()
	serverId := uuid.New().String()
	createLogger := w.logger.With(
		"server_id", serverId,
	)
	createLogger.Info("Received create server message", "Value", string(message.Value))

	serverConfig := &types.CreateServerConfig{}
	if err := json.Unmarshal(message.Value, serverConfig); err != nil {
		createLogger.Error("Failed to unmarshal server config", "error", err)
		return err
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
		//todo: handle failure
		return err
	}

	createLogger.Info("server created successfully")
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

	w.logger.Info("Worker started and listening for messages")
	redpandaConsumer.ReadMessage(w.handleCreateServer)
	w.logger.Info("Closing worker")

	return nil
}
