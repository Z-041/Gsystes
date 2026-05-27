package migration

import (
	"github.com/gsystes/backend/internal/data/model"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
	)
	if err != nil {
		logger.Error("auto migration failed", logger.ErrorField(err))
		return err
	}
	logger.Info("database auto migration completed")
	return nil
}
