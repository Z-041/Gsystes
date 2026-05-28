package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/handler"
	mid "github.com/gsystes/backend/internal/communication/middleware"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	permHandler *handler.PermissionHandler,
	operationLogHandler *handler.OperationLogHandler,
	operationLogMid *mid.OperationLogMiddleware,
) *gin.Engine {
	r := gin.New()

	r.Use(mid.Recovery())
	r.Use(mid.CORS())
	r.Use(mid.RequestLogger())

	api := r.Group("/api/v1")
	api.Use(operationLogMid.Handle())
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.GET("/menus", mid.AuthRequired(), userHandler.GetCurrentMenus)
			auth.GET("/permissions", mid.AuthRequired(), userHandler.GetCurrentPermissions)
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
			users.PUT("/:id/role", userHandler.AssignRole)
			users.POST("/batch/role", userHandler.BatchAssignRole)
			users.GET("/by-role/:roleId", userHandler.GetUsersByRole)
		}

		roles := api.Group("/roles")
		roles.Use(mid.AuthRequired())
		{
			roles.GET("/all", roleHandler.ListAll)
			roles.POST("", roleHandler.Create)
			roles.PUT("/:id", roleHandler.Update)
			roles.DELETE("/:id", roleHandler.Delete)
			roles.GET("/:id", roleHandler.Get)
			roles.GET("", roleHandler.List)
			roles.POST("/:id/permissions", roleHandler.AssignPermissions)
			roles.GET("/:id/permissions", roleHandler.GetPermissions)
		}

		permissions := api.Group("/permissions")
		permissions.Use(mid.AuthRequired())
		{
			permissions.GET("/all", permHandler.ListAll)
			permissions.POST("", permHandler.Create)
			permissions.PUT("/:id", permHandler.Update)
			permissions.DELETE("/:id", permHandler.Delete)
			permissions.GET("/:id", permHandler.Get)
			permissions.GET("", permHandler.List)
		}

		logs := api.Group("/logs")
		logs.Use(mid.AuthRequired())
		{
			logs.GET("", operationLogHandler.List)
		}
	}

	return r
}
