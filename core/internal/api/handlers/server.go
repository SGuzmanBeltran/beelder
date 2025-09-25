package handlers

import (
	"beelder/internal/api/services"
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

	servers.Post("", validation.ValidateBody[CreateServerConfig], h.createServer)
}

type CreateServerConfig struct {
	ServerType      string `json:"server_type" validate:"required"`
	QuantityMembers int    `json:"quantity_members" validate:"required,min=1"`
	PlanType        string `json:"plan_type" validate:"required"`
}

func (h *ServerHandler) createServer(c *fiber.Ctx) error {
	var serverConfig CreateServerConfig

	if err := c.BodyParser(serverConfig); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Cannot parse request body",
        })
    }

	result := h.serverService.CreateServer()
	if result != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": result.Error(),
		})
	}
	return c.SendStatus(fiber.StatusCreated)
}
