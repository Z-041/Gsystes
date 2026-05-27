package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gsystes/backend/internal/infrastructure/logger"
)

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        c.Next()

        latency := time.Since(start)
        statusCode := c.Writer.Status()

        logger.Info("request",
            logger.StringField("method", method),
            logger.StringField("path", path),
            logger.IntField("status", statusCode),
            logger.DurationField("latency", latency),
            logger.StringField("client_ip", c.ClientIP()),
        )
    }
}