package repository

import "github.com/gsystes/backend/internal/domain/entity"

type OperationLogRepository interface {
	Create(log *entity.OperationLog) error
	FindByPage(page, pageSize int) ([]entity.OperationLog, int64, error)
}
