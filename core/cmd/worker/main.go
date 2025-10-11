package main

import (
	"beelder/internal/worker"
	"beelder/internal/worker/builder"
	"log/slog"
)

func main() {
	logger := slog.Default()
	builder := builder.NewBuilder()
	worker := worker.NewWorker(builder)

	if err := worker.Start(); err != nil {
		logger.Error("worker failed", "error", err)
	}

	logger.Info("worker exited")
}
