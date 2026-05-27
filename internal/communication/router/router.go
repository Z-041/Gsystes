package router

import (
	"github.com/gin-gonic/gin"
	mid "github.com/gsystes/backend/internal/communication/middleware"
	"github.com/gsystes/backend/internal/communication/handler"
)

func SetupRouter(
	userHandler *handler.UserHandler,
) *gin.Engine {
	r := gin.New()

	r.Use(mid.Recovery())
	r.Use(mid.CORS())
	r.Use(mid.RequestLogger())

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
		}

		users := api.Group("/users")
		users.Use(mid.AuthRequired())
		{
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
			users.GET("/:id", userHandler.Get)
			users.GET("", userHandler.List)
			users.PUT("/password", userHandler.ChangePassword)
		}
	}

	return r
}