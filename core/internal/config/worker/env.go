package config

import (
	"beelder/internal/config"
	"encoding/json"
	"log/slog"
	"os"
)

type BuilderConfig struct {
	MaxConcurrentBuilds int32 `json:"max_concurrent_builds"`
	MaxAliveServers    	int32 `json:"max_alive_servers"`
	BuildTimeout        int32 `json:"timeout_seconds"`
}

type WorkerConfig struct {
	Broker     string
	ConsumerTopic string
	ProducerTopic  string
	GroupID string
	DockerHost string
	BuilderConfig BuilderConfig
}

var WorkerEnvs = initConfig()

func initConfig() WorkerConfig {
	configLogger := slog.Default().With("component", "config")
	configLogger.Info("Loading worker environment variables")
	if err := config.LoadEnv("worker"); err != nil {
		configLogger.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	var builderConfig BuilderConfig
	builderConfigString := config.GetEnv("BUILDER_CONFIG")
	err := json.Unmarshal([]byte(builderConfigString), &builderConfig); if err != nil {
		configLogger.Error("Error parsing BUILDER_CONFIG", "error", err)
		os.Exit(1)
	}

	config := WorkerConfig{
		Broker:     config.GetEnv("BROKER"),
		ConsumerTopic: config.GetEnv("CONSUMER_TOPIC"),
		ProducerTopic:  config.GetEnv("PRODUCER_TOPIC"),
		GroupID: config.GetEnv("GROUP_ID"),
		DockerHost: config.GetEnv("DOCKER_HOST"),
		BuilderConfig: builderConfig,
	}

	configLogger.Info("Worker environment variables loaded successfully!")

	return config
}