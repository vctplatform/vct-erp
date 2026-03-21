package realtime

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// WebsocketHandler upgrades executive finance requests into websocket sessions.
type WebsocketHandler struct {
	hub          *Hub
	lifecycle    context.Context
	upgrader     websocket.Upgrader
	allowedRoles map[string]struct{}
	roleHeader   string
	actorHeader  string
}

// NewWebsocketHandler constructs the finance websocket HTTP adapter.
func NewWebsocketHandler(lifecycle context.Context, hub *Hub, allowedOrigins []string, roleHeader string, actorHeader string, allowedRoles ...string) *WebsocketHandler {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		trimmed := strings.ToLower(strings.TrimSpace(role))
		if trimmed != "" {
			roleSet[trimmed] = struct{}{}
		}
	}

	return &WebsocketHandler{
		hub:          hub,
		lifecycle:    lifecycle,
		allowedRoles: roleSet,
		roleHeader:   fallbackHeader(roleHeader, "X-App-Role"),
		actorHeader:  fallbackHeader(actorHeader, "X-Actor-ID"),
		upgrader: websocket.Upgrader{
			CheckOrigin: originChecker(allowedOrigins),
		},
	}
}

// ServeHTTP upgrades the request and starts the websocket client loops.
func (h *WebsocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	role := strings.ToLower(firstNonEmpty(
		strings.TrimSpace(r.Header.Get(h.roleHeader)),
		strings.TrimSpace(r.URL.Query().Get("role")),
	))
	actorID := firstNonEmpty(
		strings.TrimSpace(r.Header.Get(h.actorHeader)),
		strings.TrimSpace(r.URL.Query().Get("actor_id")),
	)
	if actorID == "" || !h.isAllowed(role) {
		writeWSJSON(w, http.StatusForbidden, map[string]string{
			"error":   "forbidden",
			"message": "finance websocket requires cfo, ceo, or system_admin credentials",
		})
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := NewClient(conn, h.hub, actorID, role)
	if !h.hub.Register(client) {
		_ = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "role not allowed"), time.Now().Add(writeWait))
		_ = conn.Close()
		return
	}

	lifecycle := h.lifecycle
	if lifecycle == nil {
		lifecycle = r.Context()
	}
	client.Serve(lifecycle)
}

func (h *WebsocketHandler) isAllowed(role string) bool {
	if len(h.allowedRoles) == 0 {
		return true
	}
	_, ok := h.allowedRoles[strings.ToLower(strings.TrimSpace(role))]
	return ok
}

func originChecker(allowedOrigins []string) func(*http.Request) bool {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			allowed[trimmed] = struct{}{}
		}
	}

	return func(r *http.Request) bool {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin == "" || len(allowed) == 0 {
			return true
		}
		if _, ok := allowed["*"]; ok {
			return true
		}
		_, ok := allowed[origin]
		return ok
	}
}

func fallbackHeader(header string, fallback string) string {
	if strings.TrimSpace(header) == "" {
		return fallback
	}
	return header
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func writeWSJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
