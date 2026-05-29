package handler

import (
	"github.com/gin-gonic/gin"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

type DashboardHandler struct {
	userRepo domainRepo.UserRepository
	roleRepo domainRepo.RoleRepository
	logRepo  domainRepo.OperationLogRepository
}

func NewDashboardHandler(
	userRepo domainRepo.UserRepository,
	roleRepo domainRepo.RoleRepository,
	logRepo domainRepo.OperationLogRepository,
) *DashboardHandler {
	return &DashboardHandler{
		userRepo: userRepo,
		roleRepo: roleRepo,
		logRepo:  logRepo,
	}
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	userCount, err := h.userRepo.Count()
	if err != nil {
		utils.InternalError(c, "failed to get user count")
		return
	}

	roleCount, err := h.roleRepo.Count()
	if err != nil {
		utils.InternalError(c, "failed to get role count")
		return
	}

	todayLogCount, _ := h.logRepo.CountToday()

	utils.Success(c, gin.H{
		"user_count":       userCount,
		"role_count":       roleCount,
		"today_log_count":  todayLogCount,
	})
}
