package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gsystes/backend/internal/communication/handler"
	mid "github.com/gsystes/backend/internal/communication/middleware"
	"github.com/gsystes/backend/internal/communication/router"
	"github.com/gsystes/backend/internal/data/migration"
	dataRepo "github.com/gsystes/backend/internal/data/repository"
	"github.com/gsystes/backend/internal/data/seed"
	domainService "github.com/gsystes/backend/internal/domain/service"
	"github.com/gsystes/backend/internal/infrastructure/async"
	"github.com/gsystes/backend/internal/infrastructure/cache"
	"github.com/gsystes/backend/internal/infrastructure/config"
	"github.com/gsystes/backend/internal/infrastructure/database"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
	"golang.org/x/sync/errgroup"

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

	eg, _ := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return database.InitDatabase(cfg.Database)
	})
	eg.Go(func() error {
		return cache.InitRedis(cfg.Redis)
	})

	if err := eg.Wait(); err != nil {
		logger.Fatal("failed to initialize infrastructure", logger.ErrorField(err))
	}

	if err := migration.AutoMigrate(database.GetDB()); err != nil {
		logger.Fatal("failed to auto migrate", logger.ErrorField(err))
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

	logWriter := async.NewOperationLogWriter(operationLogRepo, 4, 4096)
	logWriter.Start()

	userOrchestration := orchestration.NewUserOrchestration(userDomainService, userRepo, roleRepo)
	roleOrchestration := orchestration.NewRoleOrchestration(roleRepo, permRepo)
	permOrchestration := orchestration.NewPermissionOrchestration(permRepo)
	logOrchestration := orchestration.NewOperationLogOrchestration(operationLogRepo)

	userHandler := handler.NewUserHandler(userOrchestration)
	roleHandler := handler.NewRoleHandler(roleOrchestration)
	permHandler := handler.NewPermissionHandler(permOrchestration)
	logHandler := handler.NewOperationLogHandler(logOrchestration)

	operationLogMid := mid.NewOperationLogMiddleware(logWriter)
	permMid := mid.NewPermissionMiddleware(roleRepo)

	r := router.SetupRouter(userHandler, roleHandler, permHandler, logHandler, operationLogMid, permMid)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		logger.Info("server starting", logger.StringField("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed to start", logger.ErrorField(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", logger.ErrorField(err))
	}

	logWriter.Stop()
	mid.StopMemoryLimiter()

	if err := cache.Close(); err != nil {
		logger.Error("failed to close redis", logger.ErrorField(err))
	}
	if err := database.Close(); err != nil {
		logger.Error("failed to close database", logger.ErrorField(err))
	}
	if err := logger.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to sync logger: %v\n", err)
	}

	logger.Info("server exited gracefully")
}
