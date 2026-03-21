package realtime

import (
	"encoding/json"
	"strings"
	"sync"

	analyticsdomain "vct-platform/backend/internal/modules/analytics/domain"
)

// Hub manages live finance websocket subscribers.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]struct{}
	roles   map[string]struct{}
}

// NewHub constructs an in-memory websocket hub for executive finance events.
func NewHub(allowedRoles ...string) *Hub {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		trimmed := strings.ToLower(strings.TrimSpace(role))
		if trimmed != "" {
			roleSet[trimmed] = struct{}{}
		}
	}

	return &Hub{
		clients: make(map[*Client]struct{}),
		roles:   roleSet,
	}
}

// Register stores a client if its role is permitted for finance telemetry.
func (h *Hub) Register(client *Client) bool {
	if h == nil || client == nil {
		return false
	}
	if len(h.roles) > 0 {
		if _, ok := h.roles[strings.ToLower(client.Role())]; !ok {
			return false
		}
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = struct{}{}
	return true
}

// Unregister removes a client from the hub.
func (h *Hub) Unregister(client *Client) {
	if h == nil || client == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, client)
}

// Broadcast sends a realtime finance event to all online executive clients.
func (h *Hub) Broadcast(event analyticsdomain.DashboardEvent) {
	if h == nil {
		return
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		select {
		case client.send <- payload:
		default:
			go client.Close()
		}
	}
}
