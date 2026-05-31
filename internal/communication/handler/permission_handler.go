package handler

import (
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

// Create godoc
// @Summary      创建权限
// @Description  创建一个新权限（菜单或按钮）
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreatePermissionRequest  true  "权限信息"
// @Success      200  {object}  utils.Response{data=object{id=uint}}
// @Failure      400  {object}  utils.Response
// @Router       /permissions [post]
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

// Update godoc
// @Summary      更新权限
// @Description  更新指定权限的信息
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  uint                          true  "权限 ID"
// @Param        body  body  dto.UpdatePermissionRequest  true  "更新信息"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /permissions/{id} [put]
func (h *PermissionHandler) Update(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var req dto.UpdatePermissionRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.permOrchestration.UpdatePermission(&orchestration.UpdatePermissionRequest{
		ID:       id,
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

// Delete godoc
// @Summary      删除权限
// @Description  删除指定权限
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      uint  true  "权限 ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /permissions/{id} [delete]
func (h *PermissionHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	if err := h.permOrchestration.DeletePermission(id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// Get godoc
// @Summary      获取权限详情
// @Description  根据 ID 获取权限详细信息
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      uint  true  "权限 ID"
// @Success      200  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /permissions/{id} [get]
func (h *PermissionHandler) Get(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	p, err := h.permOrchestration.GetPermission(id)
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

// List godoc
// @Summary      权限列表
// @Description  分页查询权限列表
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page       query  int  false  "页码（默认 1）"
// @Param        page_size  query  int  false  "每页条数（默认 10，最大 100）"
// @Success      200  {object}  utils.PageResult
// @Failure      500  {object}  utils.Response
// @Router       /permissions [get]
func (h *PermissionHandler) List(c *gin.Context) {
	pg := utils.GetPagination(c)

	perms, total, err := h.permOrchestration.ListPermissions(pg.Page, pg.PageSize)
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

	utils.PageSuccess(c, permList, total, pg.Page, pg.PageSize)
}

// ListAll godoc
// @Summary      获取所有权限（精简）
// @Description  返回所有权限列表，不包含分页
// @Tags         权限管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Router       /permissions/all [get]
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
