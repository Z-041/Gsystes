package websocket

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gsystes/backend/internal/infrastructure/auth"
	"github.com/gsystes/backend/internal/infrastructure/logger"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleUpgrade(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.URL.Query().Get("token")
		if tokenStr == "" {
			http.Error(w, "token is required", http.StatusUnauthorized)
			return
		}

		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Error("websocket upgrade failed", logger.ErrorField(err))
			return
		}

		client := NewClient(hub, conn, claims.UserID, claims.Username)
		hub.register <- client

		connectedPayload := ConnectedPayload{
			ClientID: client.id,
			Username: claims.Username,
			Message:  "connected",
		}
		payloadBytes, _ := json.Marshal(connectedPayload)

		msg := Message{
			Type:      TypeConnected,
			Payload:   payloadBytes,
			Timestamp: time.Now().UnixMilli(),
		}
		raw, _ := json.Marshal(msg)
		select {
		case client.send <- raw:
		default:
		}

		go client.WritePump()
		go client.ReadPump()
	}
}
