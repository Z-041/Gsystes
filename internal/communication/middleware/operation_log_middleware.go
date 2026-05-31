package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/websocket"
	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/infrastructure/async"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
)

type OperationLogMiddleware struct {
	writer *async.OperationLogWriter
	wsHub  *websocket.Hub
}

func NewOperationLogMiddleware(writer *async.OperationLogWriter, wsHub *websocket.Hub) *OperationLogMiddleware {
	return &OperationLogMiddleware{writer: writer, wsHub: wsHub}
}

var sensitiveKeys = []string{"password", "old_password", "new_password", "secret", "token", "access_token", "refresh_token"}

func sanitizeBody(body string) string {
	if body == "" {
		return ""
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		return body
	}
	for _, key := range sensitiveKeys {
		if _, ok := m[key]; ok {
			m[key] = "***"
		}
	}
	result, _ := json.Marshal(m)
	return string(result)
}

var skipPaths = []string{"/uploads/", "/swagger/", "/health", "/favicon.ico"}

func shouldSkip(path string) bool {
	for _, p := range skipPaths {
		if strings.Contains(path, p) {
			return true
		}
	}
	return false
}

func extractModule(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	if len(parts) == 1 {
		return parts[0]
	}
	return "other"
}

func methodToAction(method string) string {
	switch method {
	case "GET":
		return "query"
	case "POST":
		return "create"
	case "PUT":
		return "update"
	case "DELETE":
		return "delete"
	case "PATCH":
		return "update"
	default:
		return method
	}
}

func (m *OperationLogMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		query := c.Request.URL.RawQuery

		if shouldSkip(path) {
			c.Next()
			return
		}

		var body string
		if c.Request.Body != nil {
			data, _ := io.ReadAll(c.Request.Body)
			body = sanitizeBody(string(data))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(data))
		}

		c.Next()

		latency := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()

		entry := &entity.OperationLog{
			Method:     method,
			Path:       path,
			Query:      query,
			Body:       body,
			StatusCode: statusCode,
			Latency:    latency,
			ClientIP:   c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}

		var username string
		if claims := infraMiddleware.GetClaims(c); claims != nil {
			entry.UserID = claims.UserID
			entry.Username = claims.Username
			username = claims.Username
		}

		m.writer.Write(entry)

		if m.wsHub != nil && m.wsHub.ClientCount() > 0 {
			m.wsHub.BroadcastLogEntry(&websocket.LogEntryPayload{
				Username:   username,
				Module:     extractModule(path),
				Action:     methodToAction(method),
				Method:     method,
				Path:       path,
				IP:         c.ClientIP(),
				Duration:   latency,
				StatusCode: statusCode,
				CreatedAt:  time.Now().Format(time.DateTime),
			})
		}
	}
}
