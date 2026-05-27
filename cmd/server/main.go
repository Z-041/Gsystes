package main

import (
	"flag"
	"fmt"

	"github.com/gsystes/backend/internal/communication/handler"
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
)

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

	userDomainService := domainService.NewUserDomainService(userRepo)

	if err := seed.InitSeedData(db, userDomainService, userRepo, roleRepo, permRepo); err != nil {
		logger.Fatal("failed to init seed data", logger.ErrorField(err))
	}

	userOrchestration := orchestration.NewUserOrchestration(userDomainService, userRepo)
	roleOrchestration := orchestration.NewRoleOrchestration(roleRepo, permRepo)
	permOrchestration := orchestration.NewPermissionOrchestration(permRepo)

	userHandler := handler.NewUserHandler(userOrchestration)
	roleHandler := handler.NewRoleHandler(roleOrchestration)
	permHandler := handler.NewPermissionHandler(permOrchestration)

	r := router.SetupRouter(userHandler, roleHandler, permHandler)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("server starting", logger.StringField("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal("server failed to start", logger.ErrorField(err))
	}
}
