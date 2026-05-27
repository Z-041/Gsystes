package database

import (
    "time"

    "github.com/gsystes/backend/internal/infrastructure/config"
    "github.com/gsystes/backend/internal/infrastructure/logger"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    gormlogger "gorm.io/gorm/logger"
)

var globalDB *gorm.DB

func InitDatabase(cfg config.DatabaseConfig) error {
    gormConfig := &gorm.Config{
        Logger: gormlogger.Default.LogMode(gormlogger.Info),
    }

    db, err := gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
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
    logger.Info("database connected successfully")
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