package main

import (
	"beelder/internal/api/handlers"
	"beelder/internal/api/services"
	config "beelder/internal/config/api"
	"beelder/internal/helpers"
	"beelder/pkg/messaging/redpanda"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/store/go_cache/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	gocache "github.com/patrickmn/go-cache"
)

func main() {
	// ctx := context.Background()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",                                                                         //This should be changed to the valid origin
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

	producer := redpanda.NewRedpandaProducer(producerConfig)
	producer.Connect()

	// Initialize shared cache
	sharedCache := setupCache()

	// Initialize services
	versionProvider := helpers.NewVersionProvider()
	serverService := services.NewServerService(producer, versionProvider)

	sse := services.NewSSEService(consumerConfig)
	sse.Run()

	// Initialize handlers with shared cache
	serverHandler := handlers.NewServerHandler(serverService, sharedCache)
	sseHandler := handlers.NewSSEHandler(sse)

	// Register routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	serverHandler.RegisterRoutes(v1)
	sseHandler.RegisterRoutes(v1)
}

// setupCache initializes a shared cache instance for all handlers
func setupCache() *cache.Cache[any] {
	gocacheClient := gocache.New(1*time.Hour, 2*time.Hour)
	gocacheStore := go_cache.NewGoCache(gocacheClient)
	return cache.New[any](gocacheStore)
}
