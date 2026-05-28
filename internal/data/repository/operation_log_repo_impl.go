package repository

import (
	"github.com/gsystes/backend/internal/data/model"
	domainEntity "github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type operationLogRepository struct {
	db *gorm.DB
}

func NewOperationLogRepository(db *gorm.DB) domainRepo.OperationLogRepository {
	return &operationLogRepository{db: db}
}

func (r *operationLogRepository) Create(log *domainEntity.OperationLog) error {
	return r.db.Create(&model.OperationLog{
		UserID:     log.UserID,
		Username:   log.Username,
		Method:     log.Method,
		Path:       log.Path,
		Query:      log.Query,
		Body:       log.Body,
		StatusCode: log.StatusCode,
		Latency:    log.Latency,
		ClientIP:   log.ClientIP,
		UserAgent:  log.UserAgent,
	}).Error
}

func (r *operationLogRepository) FindByPage(page, pageSize int) ([]domainEntity.OperationLog, int64, error) {
	var models []model.OperationLog
	var total int64

	if err := r.db.Model(&model.OperationLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]domainEntity.OperationLog, len(models))
	for i, m := range models {
		entities[i] = domainEntity.OperationLog{
			ID:         m.ID,
			UserID:     m.UserID,
			Username:   m.Username,
			Method:     m.Method,
			Path:       m.Path,
			Query:      m.Query,
			Body:       m.Body,
			StatusCode: m.StatusCode,
			Latency:    m.Latency,
			ClientIP:   m.ClientIP,
			UserAgent:  m.UserAgent,
			CreatedAt:  m.CreatedAt,
		}
	}
	return entities, total, nil
}