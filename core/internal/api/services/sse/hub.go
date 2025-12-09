package sse

import (
	"sync"

	"github.com/google/uuid"
)

type ProgressEvent struct {
	ServerID string `json:"server_id"`
	Status string `json:"status"`
	Stage string `json:"stage"`
	Message string `json:"message"`
}

type Client struct {
	ID string
	ServerID string
	Channel chan ProgressEvent
}

type Hub struct {
	clients map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan ProgressEvent
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan ProgressEvent, 100),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.ID] = client
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.Channel)
			}
			h.mu.Unlock()
		case event := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				if h.shouldSendToClient(client, event) {
					client.Channel <- event
				}
			}
			h.mu.RUnlock()
		}
	}
}

// shouldSendToClient determines if an event should be sent to a client
func (h *Hub) shouldSendToClient(client *Client, event ProgressEvent) bool {
	// If client is watching a specific server, must match
	if client.ServerID != "" && client.ServerID != event.ServerID {
		return false
	}

	return true
}

// RegisterClient registers a new client
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient unregisters a client
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastEvent broadcasts an event to all matching clients
func (h *Hub) BroadcastEvent(event ProgressEvent) {
	h.broadcast <- event
}

// NewClient creates a new client instance
func NewClient(serverID string) *Client {
	return &Client{
		ID:       uuid.New().String(),
		ServerID: serverID,
		Channel:  make(chan ProgressEvent, 20), // Buffer size for event bursts
	}
}



