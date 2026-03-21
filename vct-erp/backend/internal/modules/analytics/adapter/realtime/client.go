package realtime

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

// Client models one websocket executive session.
type Client struct {
	conn    *websocket.Conn
	hub     *Hub
	actorID string
	role    string
	send    chan []byte
}

// NewClient constructs a websocket client wrapper.
func NewClient(conn *websocket.Conn, hub *Hub, actorID string, role string) *Client {
	return &Client{
		conn:    conn,
		hub:     hub,
		actorID: strings.TrimSpace(actorID),
		role:    strings.ToLower(strings.TrimSpace(role)),
		send:    make(chan []byte, 32),
	}
}

// Role returns the normalized application role attached to the client.
func (c *Client) Role() string {
	if c == nil {
		return ""
	}
	return c.role
}

// ActorID returns the actor identifier associated with the session.
func (c *Client) ActorID() string {
	if c == nil {
		return ""
	}
	return c.actorID
}

// Serve starts the read and write loops until the context or websocket closes.
func (c *Client) Serve(ctx context.Context) {
	if c == nil || c.conn == nil {
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go c.writePump(ctx)
	c.readPump(ctx)
}

// Close terminates the websocket session and unregisters the client from the hub.
func (c *Client) Close() {
	if c == nil {
		return
	}
	if c.hub != nil {
		c.hub.Unregister(c)
	}
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

func (c *Client) readPump(ctx context.Context) {
	defer c.Close()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if _, _, err := c.conn.ReadMessage(); err != nil {
				return
			}
		}
	}
}

func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			_ = c.writeControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutdown"))
			return
		case payload, ok := <-c.send:
			if !ok {
				_ = c.writeControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "hub closed"))
				return
			}
			if err := c.writeMessage(websocket.TextMessage, payload); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.writeControl(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) writeMessage(messageType int, payload []byte) error {
	if c == nil || c.conn == nil {
		return http.ErrServerClosed
	}
	_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(messageType, payload)
}

func (c *Client) writeControl(messageType int, payload []byte) error {
	if c == nil || c.conn == nil {
		return http.ErrServerClosed
	}
	return c.conn.WriteControl(messageType, payload, time.Now().Add(writeWait))
}
