package main

import (
	"flag"
	"fmt"

	"github.com/gsystes/backend/internal/communication/handler"
	mid "github.com/gsystes/backend/internal/communication/middleware"
	"github.com/gsystes/backend/internal/communication/router"
	"github.com/gsystes/backend/internal/data/migration"
	dataRepo "github.com/gsystes/backend/internal/data/repository"
	"github.com/gsystes/backend/internal/data/seed"
	domainService "github.com/gsystes/backend/internal/domain/service"
	"github.com/gsystes/backend/internal/infrastructure/cache"
	"github.com/gsystes/backend/internal/infrastructure/config"
	"github.com/gsystes/backend/internal/infrastructure/database"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"

	_ "github.com/gsystes/backend/docs/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Gsystes Backend API
// @version         1.0
// @description     中后台管理系统后端 API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and JWT token.
func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	if err := logger.InitLogger(cfg.Log); err != nil {
		panic(fmt.Sprintf("failed to init logger: %v", err))
	}

	if err := database.InitDatabase(cfg.Database); err != nil {
		logger.Fatal("failed to init database", logger.ErrorField(err))
	}

	if err := migration.AutoMigrate(database.GetDB()); err != nil {
		logger.Fatal("failed to auto migrate", logger.ErrorField(err))
	}

	if err := cache.InitRedis(cfg.Redis); err != nil {
		logger.Warn("redis init failed, continuing without cache", logger.ErrorField(err))
	}

	db := database.GetDB()
	userRepo := dataRepo.NewUserRepository(db)
	roleRepo := dataRepo.NewRoleRepository(db)
	permRepo := dataRepo.NewPermissionRepository(db)
	operationLogRepo := dataRepo.NewOperationLogRepository(db)

	userDomainService := domainService.NewUserDomainService(userRepo)

	if err := seed.InitSeedData(db, userDomainService, userRepo, roleRepo, permRepo); err != nil {
		logger.Fatal("failed to init seed data", logger.ErrorField(err))
	}

	userOrchestration := orchestration.NewUserOrchestration(userDomainService, userRepo, roleRepo)
	roleOrchestration := orchestration.NewRoleOrchestration(roleRepo, permRepo)
	permOrchestration := orchestration.NewPermissionOrchestration(permRepo)
	logOrchestration := orchestration.NewOperationLogOrchestration(operationLogRepo)

	userHandler := handler.NewUserHandler(userOrchestration)
	roleHandler := handler.NewRoleHandler(roleOrchestration)
	permHandler := handler.NewPermissionHandler(permOrchestration)
	logHandler := handler.NewOperationLogHandler(logOrchestration)

	operationLogMid := mid.NewOperationLogMiddleware(operationLogRepo)
	permMid := mid.NewPermissionMiddleware(roleRepo)

	r := router.SetupRouter(userHandler, roleHandler, permHandler, logHandler, operationLogMid, permMid)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("server starting", logger.StringField("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal("server failed to start", logger.ErrorField(err))
	}
}
