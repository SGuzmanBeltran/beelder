package main

import (
	"bee-builder/internal/handlers"
	"bee-builder/internal/services"
	"bee-builder/pkg/messaging/redpanda"
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
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization", // Added Authorization
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE",
		AllowCredentials: true,
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
		Brokers: []string{"localhost:29092"},
		Topic:   "test",
	}


	// Initializa services
	serverService := services.NewServerService(producerConfig)

	// Initialize handlers
	serverHandler := handlers.NewServerHandler(serverService)

	// Register routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	serverHandler.RegisterRoutes(v1)
}