package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gsystes/backend/internal/infrastructure/logger"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			logger.Info("websocket client connected",
				logger.StringField("client_id", client.id),
				logger.StringField("username", client.username),
			)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			logger.Info("websocket client disconnected",
				logger.StringField("client_id", client.id),
			)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastJSON(msgType MessageType, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		logger.Error("failed to marshal websocket payload", logger.ErrorField(err))
		return
	}

	msg := Message{
		Type:      msgType,
		Payload:   data,
		Timestamp: time.Now().UnixMilli(),
	}

	raw, err := json.Marshal(msg)
	if err != nil {
		logger.Error("failed to marshal websocket message", logger.ErrorField(err))
		return
	}

	select {
	case h.broadcast <- raw:
	default:
		logger.Warn("websocket broadcast channel full, dropping message")
	}
}

func (h *Hub) BroadcastLogEntry(payload *LogEntryPayload) {
	h.BroadcastJSON(TypeLogEntry, payload)
}

func (h *Hub) BroadcastStatUpdate(payload *StatUpdatePayload) {
	h.BroadcastJSON(TypeStatUpdate, payload)
}

func (h *Hub) BroadcastNotification(payload *NotificationPayload) {
	h.BroadcastJSON(TypeNotification, payload)
}

func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
