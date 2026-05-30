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

	"github.com/gsystes/backend/internal/data/migration"
	dataRepo "github.com/gsystes/backend/internal/data/repository"
	"github.com/gsystes/backend/internal/data/seed"
	"github.com/gsystes/backend/internal/domain/repository"
	domainService "github.com/gsystes/backend/internal/domain/service"
	"github.com/gsystes/backend/internal/infrastructure/cache"
	"github.com/gsystes/backend/internal/infrastructure/config"
	"github.com/gsystes/backend/internal/infrastructure/database"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	_ "github.com/gsystes/backend/docs/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type appRepos struct {
	userRepo         repository.UserRepository
	roleRepo         repository.RoleRepository
	permRepo         repository.PermissionRepository
	operationLogRepo repository.OperationLogRepository
}

func initRepos(db *gorm.DB) *appRepos {
	return &appRepos{
		userRepo:         dataRepo.NewUserRepository(db),
		roleRepo:         dataRepo.NewRoleRepository(db),
		permRepo:         dataRepo.NewPermissionRepository(db),
		operationLogRepo: dataRepo.NewOperationLogRepository(db),
	}
}

func main() {
	_ = godotenv.Load()

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

	db := database.GetDB()
	if err := migration.AutoMigrate(db); err != nil {
		logger.Fatal("failed to auto migrate", logger.ErrorField(err))
	}

	repos := initRepos(db)
	userDomainService := domainService.NewUserDomainService(repos.userRepo)

	if err := seed.InitSeedData(db, userDomainService, repos.userRepo, repos.roleRepo, repos.permRepo); err != nil {
		logger.Fatal("failed to init seed data", logger.ErrorField(err))
	}

	app := SetupContainer(repos)

	r := app.Engine
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

	app.Shutdown()

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
