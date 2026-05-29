package repository

import (
	"time"

	"github.com/gsystes/backend/internal/domain/entity"
)

type LogFilter struct {
	Username  string
	Method    string
	Path      string
	StartTime *time.Time
	EndTime   *time.Time
}

type OperationLogRepository interface {
	Create(log *entity.OperationLog) error
	FindByPage(page, pageSize int, filter *LogFilter) ([]entity.OperationLog, int64, error)
	CountToday() (int64, error)
}
