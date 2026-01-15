package handlers

import (
	"beelder/internal/api/services"
	"beelder/internal/types"
	"beelder/pkg/validation"
	"context"
	"log/slog"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/gofiber/fiber/v2"
)

type ServerHandler struct {
	serverService *services.ServerService
	cache         *cache.Cache[any]
}

func NewServerHandler(serverService *services.ServerService, cache *cache.Cache[any]) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
		cache:         cache,
	}
}

func (h *ServerHandler) RegisterRoutes(routes fiber.Router) {
	servers := routes.Group("/server")

	servers.Post("", validation.ValidateBody[types.CreateServerConfig], h.createServer)
	servers.Get("/recommended-plans", validation.ValidateQuery[types.RecommendationServerParams], h.getRecommendedPlans)
	servers.Get("/:server_type/versions", h.getServerVersions)
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
		"name":    serverConfig.Name,
		"id":      serverId,
	})
}

func (h *ServerHandler) getRecommendedPlans(c *fiber.Ctx) error {
	// Get the validated params from context
	params := c.Locals("validated").(*types.RecommendationServerParams)

	plans, err := h.serverService.GetRecommendedPlans(params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": plans,
		})
}

func (h *ServerHandler) getServerVersions(c *fiber.Ctx) error {
	logger := slog.Default()
	serverType := c.Params("server_type")
	ctx := context.Background()

	cacheKey := "server__versions:" + serverType
	if cachedData, err := h.cache.Get(ctx, cacheKey); err == nil {
		logger.Info("Cache hit for server versions", "server_type", serverType, "cache_key", cacheKey)
		c.Set("X-Cache", "hit")
		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{
				"data": fiber.Map{
					"server_type": serverType,
					"versions":    cachedData,
				},
			})
	}

	versions, err := h.serverService.GetServerVersions(serverType)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Store in cache
	if err := h.cache.Set(ctx, cacheKey, versions); err != nil {
		logger.Error("Failed to store versions in cache", "server_type", serverType, "cache_key", cacheKey, "error", err)
	}

	c.Set("X-Cache", "miss")
	return c.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"data": fiber.Map{
				"server_type": serverType,
				"versions":    versions,
			},
		})
}
