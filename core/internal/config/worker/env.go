package config

import (
	"beelder/internal/config"
	"log"
)

type WorkerConfig struct {
	Broker     string
	ConsumerTopic string
	ProducerTopic  string
	GroupID string
	DockerHost string
}

var WorkerEnvs = initConfig()

func initConfig() WorkerConfig {
	if err := config.LoadEnv("worker"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	config := WorkerConfig{
		Broker:     config.GetEnv("BROKER"),
		ConsumerTopic: config.GetEnv("CONSUMER_TOPIC"),
		ProducerTopic:  config.GetEnv("PRODUCER_TOPIC"),
		GroupID: config.GetEnv("GROUP_ID"),
	}

	return config
}