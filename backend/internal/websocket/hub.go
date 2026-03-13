package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID string
	Role   string
	Conn   *websocket.Conn
	Send   chan []byte
}

type Message struct {
	ID        string      `json:"id,omitempty"`
	Type      string      `json:"type"`
	Title     string      `json:"title"`
	Body      string      `json:"body"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (h *Hub) Register() chan *Client {
	return h.register
}

func (h *Hub) Unregister() chan *Client {
	return h.unregister
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// If client already exists, close old connection
			if old, exists := h.clients[client.UserID]; exists {
				old.Conn.Close()
				close(old.Send)
			}
			h.clients[client.UserID] = client
			h.mu.Unlock()
			fmt.Printf("User %s connected via WebSocket\n", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
				fmt.Printf("User %s disconnected from WebSocket\n", client.UserID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			// Optional: System-wide broadcast logic
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.Send <- h.mustMarshal(message):
				default:
					// If send is blocked, drop message/client
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) mustMarshal(data interface{}) []byte {
	b, _ := json.Marshal(data)
	return b
}

func (h *Hub) SendToUser(userID string, msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[userID]; ok {
		if msg.Timestamp.IsZero() {
			msg.Timestamp = time.Now()
		}
		select {
		case client.Send <- h.mustMarshal(msg):
		default:
			// Skip if busy
		}
	}
}

func (h *Hub) BroadcastToOwners(msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	payload := h.mustMarshal(msg)
	for _, client := range h.clients {
		if client.Role == "owner" {
			select {
			case client.Send <- payload:
			default:
			}
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump handles inbound messages (mainly for heartbeats/connection management)
func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// We don't expect messages from client for now, just keep-alive
	}
}
