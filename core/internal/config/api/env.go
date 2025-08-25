package config

import (
	"beelder/internal/config"
	"log"
)

type ApiConfig struct {
	ServerCommdansTopic string
	Broker              string
}

var ApiEnvs = initConfig()

func initConfig() ApiConfig {
	if err := config.LoadEnv("api"); err != nil {
		log.Fatal("Error loading .env file")
	}

	config := ApiConfig{
		ServerCommdansTopic: config.GetEnv("SERVER_COMMANDS_TOPIC"),
		Broker:              config.GetEnv("BROKER"),
	}

	return config
}
