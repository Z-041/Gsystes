package websocket

import "encoding/json"

type MessageType string

const (
	TypeLogEntry     MessageType = "log_entry"
	TypeStatUpdate   MessageType = "stat_update"
	TypeConnected    MessageType = "connected"
	TypeNotification MessageType = "notification"
	TypePong         MessageType = "pong"
)

type Message struct {
	Type      MessageType     `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp int64           `json:"timestamp"`
}

type LogEntryPayload struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Module     string `json:"module"`
	Action     string `json:"action"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	IP         string `json:"ip"`
	Duration   int64  `json:"duration"`
	StatusCode int    `json:"status_code"`
	CreatedAt  string `json:"created_at"`
}

type StatUpdatePayload struct {
	UserCount     int64 `json:"user_count"`
	RoleCount     int64 `json:"role_count"`
	TodayLogCount int64 `json:"today_log_count"`
}

type ConnectedPayload struct {
	ClientID string `json:"client_id"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type NotificationPayload struct {
	Username string `json:"username"`
	Title    string `json:"title"`
	Message  string `json:"message"`
}
