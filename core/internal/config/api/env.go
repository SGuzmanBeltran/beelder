package config

import (
	"beelder/internal/config"
	"log"
)

type ApiConfig struct {
	ServerCommdansTopic string
	ServerProgressTopic string
	GroupID             string
	Broker              string
}

var ApiEnvs = initConfig()

func initConfig() ApiConfig {
	if err := config.LoadEnv("api"); err != nil {
		log.Fatal("Error loading .env file")
	}

	config := ApiConfig{
		ServerCommdansTopic: config.GetEnv("SERVER_COMMANDS_TOPIC"),
		ServerProgressTopic: config.GetEnv("SERVER_PROGRESS_TOPIC"),
		GroupID:             config.GetEnv("GROUP_ID"),
		Broker:              config.GetEnv("BROKER"),
	}

	return config
}
