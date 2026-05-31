package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gsystes/backend/internal/infrastructure/logger"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var idCounter struct {
	mu sync.Mutex
	n  uint64
}

func nextClientID() string {
	idCounter.mu.Lock()
	defer idCounter.mu.Unlock()
	idCounter.n++
	return fmt.Sprintf("ws-%d-%d", time.Now().UnixNano(), idCounter.n)
}

type Client struct {
	id         string
	username   string
	userID     uint
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	sendFailed bool
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uint, username string) *Client {
	return &Client{
		id:       nextClientID(),
		username: username,
		userID:   userID,
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 64),
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				logger.Warn("websocket read error", logger.ErrorField(err))
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		if msg.Type == TypePong {
			c.conn.SetReadDeadline(time.Now().Add(pongWait))
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
