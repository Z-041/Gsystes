package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/dto"
	"github.com/gsystes/backend/internal/infrastructure/utils"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
)

type RoleHandler struct {
	roleOrchestration *orchestration.RoleOrchestration
}

func NewRoleHandler(roleOrchestration *orchestration.RoleOrchestration) *RoleHandler {
	return &RoleHandler{roleOrchestration: roleOrchestration}
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.CreateRoleRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	role, err := h.roleOrchestration.CreateRole(&orchestration.CreateRoleRequest{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
	})
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"id": role.ID})
}

func (h *RoleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	var req dto.UpdateRoleRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.roleOrchestration.UpdateRole(&orchestration.UpdateRoleRequest{
		ID:          uint(id),
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
	}); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	if err := h.roleOrchestration.DeleteRole(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *RoleHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	role, err := h.roleOrchestration.GetRole(uint(id))
	if err != nil {
		utils.NotFound(c, "role not found")
		return
	}

	utils.Success(c, gin.H{
		"id":          role.ID,
		"name":        role.Name,
		"code":        role.Code,
		"description": role.Description,
		"status":      role.Status,
		"permissions": role.Permissions,
		"created_at":  role.CreatedAt,
		"updated_at":  role.UpdatedAt,
	})
}

func (h *RoleHandler) List(c *gin.Context) {
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

	roles, total, err := h.roleOrchestration.ListRoles(page, pageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	roleList := make([]gin.H, len(roles))
	for i, r := range roles {
		roleList[i] = gin.H{
			"id":          r.ID,
			"name":        r.Name,
			"code":        r.Code,
			"description": r.Description,
			"status":      r.Status,
			"created_at":  r.CreatedAt,
		}
	}

	utils.PageSuccess(c, roleList, total, page, pageSize)
}

func (h *RoleHandler) ListAll(c *gin.Context) {
	roles, err := h.roleOrchestration.ListAllRoles()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	roleList := make([]gin.H, len(roles))
	for i, r := range roles {
		roleList[i] = gin.H{
			"id":   r.ID,
			"name": r.Name,
			"code": r.Code,
		}
	}

	utils.Success(c, roleList)
}

func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	var req dto.AssignPermissionsRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.roleOrchestration.AssignPermissions(&orchestration.AssignPermissionsRequest{
		RoleID:        uint(id),
		PermissionIDs: req.PermissionIDs,
	}); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *RoleHandler) GetPermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	permissions, err := h.roleOrchestration.GetRolePermissions(uint(id))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, permissions)
}
