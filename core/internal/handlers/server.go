package handlers

import (
	"bee-builder/internal/services"

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
	servers := routes.Group("/servers")

	servers.Post("", h.createServer)
}

func (h *ServerHandler) createServer(c *fiber.Ctx) error {
	result := h.serverService.CreateServer()
	if result != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNotImplemented)
}
