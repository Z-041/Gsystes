package service

import (
	"github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
)

type OperationLogOrchestration struct {
	logRepo domainRepo.OperationLogRepository
}

func NewOperationLogOrchestration(logRepo domainRepo.OperationLogRepository) *OperationLogOrchestration {
	return &OperationLogOrchestration{logRepo: logRepo}
}

func (s *OperationLogOrchestration) ListLogs(page, pageSize int) ([]entity.OperationLog, int64, error) {
	return s.logRepo.FindByPage(page, pageSize)
}