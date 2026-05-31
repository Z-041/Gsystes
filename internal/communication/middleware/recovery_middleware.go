package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					logger.AnyField("error", err),
					logger.AnyField("stack", string(debug.Stack())),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Response{
					Code:    -1,
					Message: "internal server error",
				})
			}
		}()
		c.Next()
	}
}
