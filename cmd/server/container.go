package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/handler"
	mid "github.com/gsystes/backend/internal/communication/middleware"
	"github.com/gsystes/backend/internal/communication/router"
	ws "github.com/gsystes/backend/internal/communication/websocket"
	domainService "github.com/gsystes/backend/internal/domain/service"
	"github.com/gsystes/backend/internal/infrastructure/async"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
)

type AppContainer struct {
	Engine    *gin.Engine
	logWriter *async.OperationLogWriter
}

func SetupContainer(repos *appRepos) *AppContainer {
	userDomainService := domainService.NewUserDomainService(repos.userRepo)

	logWriter := async.NewOperationLogWriter(repos.operationLogRepo, 4, 4096)
	logWriter.Start()

	wsHub := ws.NewHub()
	go wsHub.Run()

	userOrchestration := orchestration.NewUserOrchestration(userDomainService, repos.userRepo, repos.roleRepo)
	roleOrchestration := orchestration.NewRoleOrchestration(repos.roleRepo, repos.permRepo)
	permOrchestration := orchestration.NewPermissionOrchestration(repos.permRepo)
	logOrchestration := orchestration.NewOperationLogOrchestration(repos.operationLogRepo)
	dashboardOrchestration := orchestration.NewDashboardOrchestration(repos.userRepo, repos.roleRepo, repos.operationLogRepo)

	eventBroadcaster := handler.NewEventBroadcaster(wsHub, repos.userRepo, repos.roleRepo, repos.operationLogRepo)

	operationLogMid := mid.NewOperationLogMiddleware(logWriter, wsHub)
	permMid := mid.NewPermissionMiddleware(repos.roleRepo)

	userHandler := handler.NewUserHandler(userOrchestration, eventBroadcaster)
	roleHandler := handler.NewRoleHandler(roleOrchestration, eventBroadcaster, permMid)
	permHandler := handler.NewPermissionHandler(permOrchestration)
	logHandler := handler.NewOperationLogHandler(logOrchestration)
	dashboardHandler := handler.NewDashboardHandler(dashboardOrchestration)

	r := router.SetupRouter(userHandler, roleHandler, permHandler, logHandler, dashboardHandler, operationLogMid, permMid, wsHub)

	return &AppContainer{
		Engine:    r,
		logWriter: logWriter,
	}
}

func (c *AppContainer) Shutdown() {
	c.logWriter.Stop()
	mid.StopMemoryLimiter()
}
