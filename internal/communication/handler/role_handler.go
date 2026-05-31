package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/dto"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
	"github.com/gsystes/backend/internal/infrastructure/utils"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
)

type RoleHandler struct {
	roleOrchestration *orchestration.RoleOrchestration
	events            *EventBroadcaster
}

func NewRoleHandler(roleOrchestration *orchestration.RoleOrchestration, events *EventBroadcaster) *RoleHandler {
	return &RoleHandler{roleOrchestration: roleOrchestration, events: events}
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

	utils.Success(c, dto.IDResponse{ID: role.ID})

	h.events.BroadcastStats()
	currentUser := infraMiddleware.GetUsername(c)
	h.events.SendNotification(currentUser, "新增角色", "角色 "+role.Name+" 已创建")
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
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var req dto.UpdateRoleRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.roleOrchestration.UpdateRole(&orchestration.UpdateRoleRequest{
		ID:          id,
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
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	if err := h.roleOrchestration.DeleteRole(id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)

	h.events.BroadcastStats()
	currentUser := infraMiddleware.GetUsername(c)
	h.events.SendNotification(currentUser, "删除角色", "角色已被移除")
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
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	role, err := h.roleOrchestration.GetRole(id)
	if err != nil {
		utils.NotFound(c, "role not found")
		return
	}

	utils.Success(c, dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Status:      role.Status,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
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
// @Success      200  {object}  utils.PageResult
// @Failure      500  {object}  utils.Response
// @Router       /roles [get]
func (h *RoleHandler) List(c *gin.Context) {
	pg := utils.GetPagination(c)

	roles, total, err := h.roleOrchestration.ListRoles(pg.Page, pg.PageSize)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	roleList := make([]dto.RoleSimpleResponse, len(roles))
	for i, r := range roles {
		roleList[i] = dto.RoleSimpleResponse{
			ID:   r.ID,
			Name: r.Name,
			Code: r.Code,
		}
	}

	utils.PageSuccess(c, roleList, total, pg.Page, pg.PageSize)
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

	roleList := make([]dto.RoleSimpleResponse, len(roles))
	for i, r := range roles {
		roleList[i] = dto.RoleSimpleResponse{
			ID:   r.ID,
			Name: r.Name,
			Code: r.Code,
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
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var req dto.AssignPermissionsRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.roleOrchestration.AssignPermissions(&orchestration.AssignPermissionsRequest{
		RoleID:        id,
		PermissionIDs: req.PermissionIDs,
	}); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)

	currentUser := infraMiddleware.GetUsername(c)
	h.events.SendNotification(currentUser, "权限变更", "角色的权限已更新")
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
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	permissions, err := h.roleOrchestration.GetRolePermissions(id)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, permissions)
}
