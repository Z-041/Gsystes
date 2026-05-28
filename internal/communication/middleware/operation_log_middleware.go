package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
)

type OperationLogMiddleware struct {
	repo domainRepo.OperationLogRepository
}

func NewOperationLogMiddleware(repo domainRepo.OperationLogRepository) *OperationLogMiddleware {
	return &OperationLogMiddleware{repo: repo}
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
			body = string(data)
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
			go func() {
				_ = m.repo.Create(entry)
			}()
		}
	}
}
