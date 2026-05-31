package service

import (
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
)

type DashboardStats struct {
	UserCount     int64
	RoleCount     int64
	TodayLogCount int64
}

type DashboardOrchestration struct {
	userRepo domainRepo.UserRepository
	roleRepo domainRepo.RoleRepository
	logRepo  domainRepo.OperationLogRepository
}

func NewDashboardOrchestration(
	userRepo domainRepo.UserRepository,
	roleRepo domainRepo.RoleRepository,
	logRepo domainRepo.OperationLogRepository,
) *DashboardOrchestration {
	return &DashboardOrchestration{
		userRepo: userRepo,
		roleRepo: roleRepo,
		logRepo:  logRepo,
	}
}

func (s *DashboardOrchestration) GetStats() (*DashboardStats, error) {
	userCount, err := s.userRepo.Count()
	if err != nil {
		return nil, err
	}

	roleCount, err := s.roleRepo.Count()
	if err != nil {
		return nil, err
	}

	todayLogCount, err := s.logRepo.CountToday()
	if err != nil {
		todayLogCount = 0
	}

	return &DashboardStats{
		UserCount:     userCount,
		RoleCount:     roleCount,
		TodayLogCount: todayLogCount,
	}, nil
}
