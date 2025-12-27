package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"app/service/dto"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 30 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512KB
)

// Client WebSocket客户端
type Client struct {
	ID       string
	UserID   int64
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	mu       sync.Mutex
	IsClosed bool
}

// Hub WebSocket连接管理器
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// Message WebSocket消息
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// NewHub 创建Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run 启动Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- marshalMessage(message):
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register 注册客户端
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister 注销客户端
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast 广播消息
func (h *Hub) Broadcast(message *Message) {
	h.broadcast <- message
}

// BroadcastLocation 广播位置更新
func (h *Hub) BroadcastLocation(entityType string, entityID string, loc *dto.LocationResp) {
	msg := &Message{
		Type: "location_update",
		Payload: mustMarshal(map[string]interface{}{
			"entity_type": entityType,
			"entity_id":   entityID,
			"location":    loc,
		}),
	}
	h.Broadcast(msg)
}

// SendToUser 发送给特定用户
func (h *Hub) SendToUser(userID int64, message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.UserID == userID {
			client.Send <- marshalMessage(message)
		}
	}
}

// NewClient 创建客户端
func NewClient(id string, userID int64, hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		ID:     id,
		UserID: userID,
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
}

// ReadPump 处理读取
func (c *Client) ReadPump(onMessage func(msg *Message) error) {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("websocket error: %v\n", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Printf("parse message error: %v\n", err)
			continue
		}

		if err := onMessage(&msg); err != nil {
			fmt.Printf("handle message error: %v\n", err)
		}
	}
}

// WritePump 处理写入
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.mu.Lock()
			if c.IsClosed {
				c.mu.Unlock()
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			c.mu.Unlock()

			if err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func marshalMessage(msg *Message) []byte {
	data, _ := json.Marshal(msg)
	return data
}

func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return data
}
