package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type ServersHandler struct {
}

func NewServersHandler() *ServersHandler {
	return &ServersHandler{}
}

func (h *ServersHandler) RegisterRoutes(routes fiber.Router) {
	servers := routes.Group("/servers")

	servers.Post("", h.createServer)
}

func (h *ServersHandler) createServer(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}