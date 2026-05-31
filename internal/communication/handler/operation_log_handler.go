package handler

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/domain/repository"
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
	pg := utils.GetPagination(c)

	var filter repository.LogFilter
	if username := c.Query("username"); username != "" {
		filter.Username = username
	}
	if method := c.Query("method"); method != "" {
		filter.Method = strings.ToUpper(method)
	}
	if path := c.Query("path"); path != "" {
		filter.Path = path
	}
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			filter.StartTime = &t
		} else if t, err := time.Parse("2006-01-02", startTime); err == nil {
			filter.StartTime = &t
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			t = t.Add(24*time.Hour - time.Second)
			filter.EndTime = &t
		} else if t, err := time.Parse("2006-01-02", endTime); err == nil {
			t = t.Add(24*time.Hour - time.Second)
			filter.EndTime = &t
		}
	}

	logs, total, err := h.logOrchestration.ListLogs(pg.Page, pg.PageSize, &filter)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	logList := make([]gin.H, len(logs))
	for i, l := range logs {
		reqBody := l.Body
		if reqBody == "" {
			reqBody = l.Query
		}
		logList[i] = gin.H{
			"id":           l.ID,
			"user_id":      l.UserID,
			"username":     l.Username,
			"method":       l.Method,
			"path":         l.Path,
			"ip":           l.ClientIP,
			"duration":     l.Latency,
			"request_body": reqBody,
			"status_code":  l.StatusCode,
			"created_at":   l.CreatedAt,
		}
	}

	utils.PageSuccess(c, logList, total, pg.Page, pg.PageSize)
}
