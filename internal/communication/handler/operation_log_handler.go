package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/infrastructure/utils"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
)

type OperationLogHandler struct {
	logOrchestration *orchestration.OperationLogOrchestration
}

func NewOperationLogHandler(logOrchestration *orchestration.OperationLogOrchestration) *OperationLogHandler {
	return &OperationLogHandler{logOrchestration: logOrchestration}
}

func (h *OperationLogHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	logs, total, err := h.logOrchestration.ListLogs(page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	logList := make([]gin.H, len(logs))
	for i, l := range logs {
		logList[i] = gin.H{
			"id":          l.ID,
			"user_id":     l.UserID,
			"username":    l.Username,
			"method":      l.Method,
			"path":        l.Path,
			"query":       l.Query,
			"status_code": l.StatusCode,
			"latency":     l.Latency,
			"client_ip":   l.ClientIP,
			"created_at":  l.CreatedAt,
		}
	}

	utils.PageSuccess(c, logList, total, page, pageSize)
}