package database

import (
	"fmt"
	"time"

	"github.com/gsystes/backend/internal/infrastructure/config"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var globalDB *gorm.DB

func openDialector(cfg config.DatabaseConfig) (gorm.Dialector, error) {
	dsn := cfg.DSN()
	switch cfg.Driver {
	case "postgres", "postgresql", "pg":
		return postgres.Open(dsn), nil
	case "mysql", "":
		return mysql.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}

func InitDatabase(cfg config.DatabaseConfig) error {
	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	}

	dialector, err := openDialector(cfg)
	if err != nil {
		logger.Error("failed to create dialector", logger.ErrorField(err))
		return err
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		logger.Error("failed to connect database", logger.ErrorField(err))
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	globalDB = db
	logger.Info("database connected successfully",
		logger.StringField("driver", cfg.Driver),
		logger.StringField("host", cfg.Host),
	)
	return nil
}

func GetDB() *gorm.DB {
	return globalDB
}

func Close() error {
	sqlDB, err := globalDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
