package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/dto"
	"github.com/gsystes/backend/internal/infrastructure/utils"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
)

type PermissionHandler struct {
	permOrchestration *orchestration.PermissionOrchestration
}

func NewPermissionHandler(permOrchestration *orchestration.PermissionOrchestration) *PermissionHandler {
	return &PermissionHandler{permOrchestration: permOrchestration}
}

func (h *PermissionHandler) Create(c *gin.Context) {
	var req dto.CreatePermissionRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	p, err := h.permOrchestration.CreatePermission(&orchestration.CreatePermissionRequest{
		Name:     req.Name,
		Code:     req.Code,
		Type:     req.Type,
		ParentID: req.ParentID,
		Path:     req.Path,
		Method:   req.Method,
		Sort:     req.Sort,
	})
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"id": p.ID})
}

func (h *PermissionHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid permission id")
		return
	}

	var req dto.UpdatePermissionRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.permOrchestration.UpdatePermission(&orchestration.UpdatePermissionRequest{
		ID:       uint(id),
		Name:     req.Name,
		Code:     req.Code,
		Type:     req.Type,
		ParentID: req.ParentID,
		Path:     req.Path,
		Method:   req.Method,
		Sort:     req.Sort,
	}); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid permission id")
		return
	}

	if err := h.permOrchestration.DeletePermission(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *PermissionHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid permission id")
		return
	}

	p, err := h.permOrchestration.GetPermission(uint(id))
	if err != nil {
		utils.NotFound(c, "permission not found")
		return
	}

	utils.Success(c, gin.H{
		"id":         p.ID,
		"name":       p.Name,
		"code":       p.Code,
		"type":       p.Type,
		"parent_id":  p.ParentID,
		"path":       p.Path,
		"method":     p.Method,
		"sort":       p.Sort,
		"created_at": p.CreatedAt,
		"updated_at": p.UpdatedAt,
	})
}

func (h *PermissionHandler) List(c *gin.Context) {
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

	perms, total, err := h.permOrchestration.ListPermissions(page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	permList := make([]gin.H, len(perms))
	for i, p := range perms {
		permList[i] = gin.H{
			"id":         p.ID,
			"name":       p.Name,
			"code":       p.Code,
			"type":       p.Type,
			"parent_id":  p.ParentID,
			"path":       p.Path,
			"method":     p.Method,
			"sort":       p.Sort,
			"created_at": p.CreatedAt,
		}
	}

	utils.PageSuccess(c, permList, total, page, pageSize)
}

func (h *PermissionHandler) ListAll(c *gin.Context) {
	perms, err := h.permOrchestration.ListAllPermissions()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	permList := make([]gin.H, len(perms))
	for i, p := range perms {
		permList[i] = gin.H{
			"id":        p.ID,
			"name":      p.Name,
			"code":      p.Code,
			"type":      p.Type,
			"parent_id": p.ParentID,
			"path":      p.Path,
			"method":    p.Method,
			"sort":      p.Sort,
		}
	}

	utils.Success(c, permList)
}
