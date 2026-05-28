package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/infrastructure/async"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
)

type OperationLogMiddleware struct {
	writer *async.OperationLogWriter
}

func NewOperationLogMiddleware(writer *async.OperationLogWriter) *OperationLogMiddleware {
	return &OperationLogMiddleware{writer: writer}
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

func (m *OperationLogMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		query := c.Request.URL.RawQuery

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

		if claims := infraMiddleware.GetClaims(c); claims != nil {
			entry.UserID = claims.UserID
			entry.Username = claims.Username
		}

		if statusCode >= 400 || method != "GET" {
			m.writer.Write(entry)
		}
	}
}
