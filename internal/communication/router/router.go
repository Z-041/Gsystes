package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/handler"
	mid "github.com/gsystes/backend/internal/communication/middleware"
	"github.com/gsystes/backend/internal/infrastructure/config"
)

func SetupRouter(
	userHandler *handler.UserHandler,
	roleHandler *handler.RoleHandler,
	permHandler *handler.PermissionHandler,
	operationLogHandler *handler.OperationLogHandler,
	dashboardHandler *handler.DashboardHandler,
	operationLogMid *mid.OperationLogMiddleware,
	permMid *mid.PermissionMiddleware,
) *gin.Engine {
	r := gin.New()

	maxBodyMB := config.GetConfig().Server.MaxBodySize
	if maxBodyMB <= 0 {
		maxBodyMB = 8
	}
	r.MaxMultipartMemory = maxBodyMB << 20

	r.Use(mid.Recovery())
	r.Use(mid.CORS())
	r.Use(mid.RequestLogger())

	r.Static("/uploads", config.GetConfig().Upload.Dir)

	rlCfg := config.GetConfig().RateLimit
	defaultLimit := mid.RateLimiter(rlCfg.DefaultRate, rlCfg.DefaultBurst, rlCfg.Window)
	loginLimit := mid.RateLimiter(rlCfg.LoginRate, rlCfg.LoginBurst, rlCfg.Window)

	api := r.Group("/api/v1")
	api.Use(operationLogMid.Handle(), defaultLimit)
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", loginLimit, userHandler.Login)
			auth.GET("/menus", mid.AuthRequired(), userHandler.GetCurrentMenus)
			auth.GET("/permissions", mid.AuthRequired(), userHandler.GetCurrentPermissions)
		}

		users := api.Group("/users")
		users.Use(mid.AuthRequired())
		{
			users.POST("", permMid.Require("user:create"), userHandler.Create)
			users.PUT("/:id", permMid.Require("user:update"), userHandler.Update)
			users.DELETE("/:id", permMid.Require("user:delete"), userHandler.Delete)
			users.GET("/:id", permMid.Require("user:read"), userHandler.Get)
			users.GET("", permMid.Require("user:read"), userHandler.List)
			users.PUT("/password", userHandler.ChangePassword)
			users.PUT("/:id/role", userHandler.AssignRole)
			users.POST("/batch/role", userHandler.BatchAssignRole)
			users.GET("/by-role/:roleId", userHandler.GetUsersByRole)

			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.POST("/avatar", userHandler.UpdateAvatar)
			users.PUT("/:id/status", userHandler.UpdateStatus)
			users.POST("/import", userHandler.ImportUsers)
			users.GET("/export", userHandler.ExportUsers)
		}

		roles := api.Group("/roles")
		roles.Use(mid.AuthRequired())
		{
			roles.GET("/all", permMid.Require("role:read"), roleHandler.ListAll)
			roles.POST("", permMid.Require("role:create"), roleHandler.Create)
			roles.PUT("/:id", permMid.Require("role:update"), roleHandler.Update)
			roles.DELETE("/:id", permMid.Require("role:delete"), roleHandler.Delete)
			roles.GET("/:id", permMid.Require("role:read"), roleHandler.Get)
			roles.GET("", permMid.Require("role:read"), roleHandler.List)
			roles.POST("/:id/permissions", permMid.Require("perm:assign"), roleHandler.AssignPermissions)
			roles.GET("/:id/permissions", permMid.Require("role:read"), roleHandler.GetPermissions)
		}

		permissions := api.Group("/permissions")
		permissions.Use(mid.AuthRequired())
		{
			permissions.GET("/all", permMid.Require("perm:manage"), permHandler.ListAll)
			permissions.POST("", permMid.Require("perm:manage"), permHandler.Create)
			permissions.PUT("/:id", permMid.Require("perm:manage"), permHandler.Update)
			permissions.DELETE("/:id", permMid.Require("perm:manage"), permHandler.Delete)
			permissions.GET("/:id", permMid.Require("perm:manage"), permHandler.Get)
			permissions.GET("", permMid.Require("perm:manage"), permHandler.List)
		}

		logs := api.Group("/logs")
		logs.Use(mid.AuthRequired())
		{
			logs.GET("", permMid.Require("log:read"), operationLogHandler.List)
		}

		dashboard := api.Group("/dashboard")
		dashboard.Use(mid.AuthRequired())
		{
			dashboard.GET("/stats", dashboardHandler.Stats)
		}
	}

	return r
}
