package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/infrastructure/config"
)

func CORS() gin.HandlerFunc {
	cfg := config.GetConfig()
	allowed := cfg.CORS.AllowedOrigins

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if len(allowed) == 0 {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			for _, o := range allowed {
				if o == origin || o == "*" {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
