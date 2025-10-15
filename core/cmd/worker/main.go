package main

import (
	"beelder/internal/worker"
	"log/slog"
)

func main() {
	logger := slog.Default()
	worker := worker.NewWorker()

	if err := worker.Start(); err != nil {
		logger.Error("worker failed", "error", err)
	}

	logger.Info("worker exited")
}
