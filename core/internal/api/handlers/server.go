package handlers

import (
	"beelder/internal/api/services"
	"beelder/internal/types"
	"beelder/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

type ServerHandler struct {
	serverService *services.ServerService
}

func NewServerHandler(serverService *services.ServerService) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
	}
}

func (h *ServerHandler) RegisterRoutes(routes fiber.Router) {
	servers := routes.Group("/server")

	servers.Post("", validation.ValidateBody[types.CreateServerConfig], h.createServer)
}

func (h *ServerHandler) createServer(c *fiber.Ctx) error {
	// Get the validated config from context
    serverConfig := c.Locals("validated").(*types.CreateServerConfig)

	serverId, err := h.serverService.CreateServer(serverConfig)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Server creation started",
        "name": serverConfig.Name,
		"id": serverId,
    })
}
