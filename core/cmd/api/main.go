package main

import (
	"beelder/internal/api/handlers"
	"beelder/internal/api/services"
	config "beelder/internal/config/api"
	"beelder/pkg/messaging/redpanda"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// ctx := context.Background()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*", //This should be changed to the valid origin
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Cache-Control, Last-Event-ID", // Added Authorization and Last-Event-ID
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE",
		AllowCredentials: false,
	}))


	setupRoutes(app)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		if err := app.Shutdown(); err != nil {
			log.Fatal("Server forced to shutdown:", err)
		}
	}()

	if err := app.Listen(":3000"); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server error:", err)
	}
}

func setupRoutes(app *fiber.App) {
	// Service configurations
	producerConfig := &redpanda.RedpandaConfig{
		Brokers: []string{config.ApiEnvs.Broker},
		Topic:   config.ApiEnvs.ServerCommdansTopic,
	}

	consumerConfig := &redpanda.RedpandaConsumerConfig{
		Brokers: []string{config.ApiEnvs.Broker},
		Topic:   config.ApiEnvs.ServerProgressTopic,
		GroupID: config.ApiEnvs.GroupID,
	}

	// Initialize services
	serverService := services.NewServerService(producerConfig)
	sse := services.NewSSEService(consumerConfig)
	sse.Run()

	// Initialize handlers
	serverHandler := handlers.NewServerHandler(serverService)
	sseHandler := handlers.NewSSEHandler(sse)

	// Register routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	serverHandler.RegisterRoutes(v1)
	sseHandler.RegisterRoutes(v1)
}