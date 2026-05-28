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

// Create godoc
// @Summary      创建角色
// @Description  创建一个新角色
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreateRoleRequest  true  "角色信息"
// @Success      200  {object}  utils.Response{data=object{id=uint}}
// @Failure      400  {object}  utils.Response
// @Router       /roles [post]
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

// Update godoc
// @Summary      更新角色
// @Description  更新指定角色的信息
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  uint                   true  "角色 ID"
// @Param        body  body  dto.UpdateRoleRequest  true  "更新信息"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /roles/{id} [put]
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

// Delete godoc
// @Summary      删除角色
// @Description  删除指定角色
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      uint  true  "角色 ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /roles/{id} [delete]
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

// Get godoc
// @Summary      获取角色详情
// @Description  根据 ID 获取角色详细信息和关联权限
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      uint  true  "角色 ID"
// @Success      200  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /roles/{id} [get]
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

// List godoc
// @Summary      角色列表
// @Description  分页查询角色列表
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page       query  int  false  "页码（默认 1）"
// @Param        page_size  query  int  false  "每页条数（默认 10，最大 100）"
// @Success      200  {object}  utils.PageResponse
// @Failure      500  {object}  utils.Response
// @Router       /roles [get]
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

// ListAll godoc
// @Summary      获取所有角色（精简）
// @Description  返回所有角色列表，只包含 ID、名称和编码
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Router       /roles/all [get]
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

// AssignPermissions godoc
// @Summary      分配角色权限
// @Description  为指定角色分配权限
// @Tags         角色-权限
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  uint                          true  "角色 ID"
// @Param        body  body  dto.AssignPermissionsRequest  true  "权限 ID 列表"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /roles/{id}/permissions [post]
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

// GetPermissions godoc
// @Summary      获取角色权限
// @Description  获取指定角色已分配的权限列表
// @Tags         角色-权限
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      uint  true  "角色 ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /roles/{id}/permissions [get]
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
