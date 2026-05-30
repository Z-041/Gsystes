package handler

import (
	"github.com/gin-gonic/gin"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

type DashboardHandler struct {
	dashboardOrchestration *orchestration.DashboardOrchestration
}

func NewDashboardHandler(dashboardOrchestration *orchestration.DashboardOrchestration) *DashboardHandler {
	return &DashboardHandler{
		dashboardOrchestration: dashboardOrchestration,
	}
}

func (h *DashboardHandler) Stats(c *gin.Context) {
	stats, err := h.dashboardOrchestration.GetStats()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"user_count":      stats.UserCount,
		"role_count":      stats.RoleCount,
		"today_log_count": stats.TodayLogCount,
	})
}
