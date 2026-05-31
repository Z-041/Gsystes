package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/dto"
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

	utils.Success(c, dto.DashboardStatsResponse{
		UserCount:     stats.UserCount,
		RoleCount:     stats.RoleCount,
		TodayLogCount: stats.TodayLogCount,
	})
}
