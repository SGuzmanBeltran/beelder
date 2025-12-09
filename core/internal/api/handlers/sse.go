package handlers

import (
	"beelder/internal/api/services"
	"beelder/internal/api/services/sse"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type SSEHandler struct {
	sseService *services.SSEService
}

func NewSSEHandler(sseService *services.SSEService) *SSEHandler {
	return &SSEHandler{
		sseService: sseService,
	}
}

func (sh *SSEHandler) RegisterRoutes(routes fiber.Router) {
	servers := routes.Group("/sse_server")

	servers.Get("/:serverID", sh.HandleSSE)
}

func (sh *SSEHandler) HandleSSE(c *fiber.Ctx) error {
	serverID := c.Params("serverID")
	if serverID == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "serverID is required",
		})
	}
	// Create new client
	client := sse.NewClient(serverID)

	log.Printf("[Handler] SSE connection attempt: ClientID=%s, ServerID=%s",
		client.ID, serverID)

	// Register client with hub
	sh.sseService.GetHub().RegisterClient(client)

	// Set SSE headers
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	log.Printf("[Handler] SSE connection established and registered: ClientID=%s, ServerID=%s",
		client.ID, serverID)

	// Stream events to client
	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		defer func() {
			sh.sseService.GetHub().UnregisterClient(client)
		}()

		// Send connection confirmation
		if err := sh.sendConnectionEvent(w, client); err != nil {
			return
		}

		// Stream provisioning events
		for event := range client.Channel {
			if err := sh.sendProgressEvent(w, event); err != nil {
				return
			}

			// Auto-close connection after completion or failure
			if event.Status == "completed" || event.Status == "failed" {
				return
			}
		}
	}))

	return nil
}

func (sh *SSEHandler) sendConnectionEvent(w *bufio.Writer, client *sse.Client) error {
	confirmData := map[string]interface{}{
		"type":      "connected",
		"clientId":  client.ID,
		"serverId":  client.ServerID,
		"message":   "Connected to server build updates",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(confirmData)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "event: connected\n")
	fmt.Fprintf(w, "data: %s\n\n", jsonData)

	return w.Flush()
}

// sendProgressEvent sends a progress event to the client
func (sh *SSEHandler) sendProgressEvent(w *bufio.Writer, event sse.ProgressEvent) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Determine event type based on status
	eventType := "progress"
	if event.Status == "completed" {
		eventType = "completed"
	} else if event.Status == "failed" {
		eventType = "failed"
	}

	fmt.Fprintf(w, "event: %s\n", eventType)
	fmt.Fprintf(w, "data: %s\n\n", jsonData)

	return w.Flush()
}
