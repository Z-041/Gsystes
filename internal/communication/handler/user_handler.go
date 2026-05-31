package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/dto"
	"github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	"github.com/gsystes/backend/internal/infrastructure/config"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
	"github.com/gsystes/backend/internal/infrastructure/utils"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
	"github.com/xuri/excelize/v2"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

type UserHandler struct {
	userOrchestration *orchestration.UserOrchestration
	events            *EventBroadcaster
}

func NewUserHandler(userOrchestration *orchestration.UserOrchestration, events *EventBroadcaster) *UserHandler {
	return &UserHandler{userOrchestration: userOrchestration, events: events}
}

// Login godoc
// @Summary      用户登录
// @Description  使用用户名和密码登录，返回 JWT Token
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        body  body  dto.LoginRequest  true  "登录信息"
// @Success      200  {object}  utils.Response{data=object{token=string,user=object{id=uint,username=string,nickname=string,avatar=string,role_id=uint}}}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Router       /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	resp, err := h.userOrchestration.Login(&orchestration.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.Success(c, dto.LoginResponse{
		Token: resp.Token,
		User: dto.UserSimple{
			ID:       resp.User.ID,
			Username: resp.User.Username,
			Nickname: resp.User.Nickname,
			Avatar:   resp.User.Avatar,
			RoleID:   resp.User.RoleID,
		},
	})
}

// Create godoc
// @Summary      创建用户
// @Description  创建一个新用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreateUserRequest  true  "用户信息"
// @Success      200  {object}  utils.Response{data=object{id=uint}}
// @Failure      400  {object}  utils.Response
// @Router       /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	user, err := h.userOrchestration.CreateUser(&orchestration.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
	})
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, dto.IDResponse{ID: user.ID})

	h.events.BroadcastStats()
	name := user.Nickname
	if name == "" {
		name = user.Username
	}
	h.events.SendNotification("系统", "新增用户", name+" 已被加入系统")
}

// Update godoc
// @Summary      更新用户
// @Description  更新指定用户的信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  uint                   true  "用户 ID"
// @Param        body  body  dto.UpdateUserRequest  true  "更新信息"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var req dto.UpdateUserRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.userOrchestration.UpdateUser(&orchestration.UpdateUserRequest{
		ID:       id,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		RoleID:   req.RoleID,
		Status:   req.Status,
	}); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// Delete godoc
// @Summary      删除用户
// @Description  删除指定用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      uint  true  "用户 ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	if err := h.userOrchestration.DeleteUser(id); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)

	h.events.BroadcastStats()
	h.events.SendNotification("系统", "删除用户", "用户已被移除")
}

// Get godoc
// @Summary      获取用户详情
// @Description  根据 ID 获取用户详细信息
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      uint  true  "用户 ID"
// @Success      200  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	user, err := h.userOrchestration.GetUser(id)
	if err != nil {
		utils.NotFound(c, "user not found")
		return
	}

	resp := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Status:    user.Status,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if user.Role != nil {
		resp.Role = &dto.RoleSimpleResponse{
			ID:   user.Role.ID,
			Name: user.Role.Name,
			Code: user.Role.Code,
		}
	}

	utils.Success(c, resp)
}

// List godoc
// @Summary      用户列表
// @Description  分页查询用户列表，支持按用户名、状态、角色筛选
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page       query  int     false  "页码（默认 1）"
// @Param        page_size  query  int     false  "每页条数（默认 10，最大 100）"
// @Param        username   query  string  false  "用户名（模糊搜索）"
// @Param        status     query  string  false  "状态（0=禁用，1=启用）"
// @Param        role_id    query  string  false  "角色 ID"
// @Success      200  {object}  utils.PageResult
// @Failure      500  {object}  utils.Response
// @Router       /users [get]
func (h *UserHandler) List(c *gin.Context) {
	pg := utils.GetPagination(c)

	filter := &domainRepo.UserFilter{
		PreloadRole: true,
	}
	if username := c.Query("username"); username != "" {
		filter.Username = username
	}
	if statusStr := c.Query("status"); statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		filter.Status = &status
	}
	if roleIDStr := c.Query("role_id"); roleIDStr != "" {
		roleID, _ := strconv.ParseUint(roleIDStr, 10, 64)
		rid := uint(roleID)
		filter.RoleID = &rid
	}

	users, total, err := h.userOrchestration.ListUsers(pg.Page, pg.PageSize, filter)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	userList := make([]dto.UserListItem, len(users))
	for i, u := range users {
		item := dto.UserListItem{
			ID:        u.ID,
			Username:  u.Username,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Phone:     u.Phone,
			Avatar:    u.Avatar,
			Status:    u.Status,
			RoleID:    u.RoleID,
			CreatedAt: u.CreatedAt,
		}
		if u.Role != nil {
			item.Role = &dto.RoleSimpleResponse{
				ID:   u.Role.ID,
				Name: u.Role.Name,
				Code: u.Role.Code,
			}
		}
		userList[i] = item
	}

	utils.PageSuccess(c, userList, total, pg.Page, pg.PageSize)
}

// ChangePassword godoc
// @Summary      修改密码
// @Description  当前用户修改自己的密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.ChangePasswordRequest  true  "密码信息"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /users/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	userID := infraMiddleware.GetUserID(c)
	if err := h.userOrchestration.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// AssignRole godoc
// @Summary      为用户分配角色
// @Description  为指定用户分配角色
// @Tags         用户-角色
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  uint            true  "用户 ID"
// @Param        body  body  object{role_id=uint}  true  "角色 ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /users/{id}/role [put]
func (h *UserHandler) AssignRole(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var req dto.AssignRoleRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.userOrchestration.AssignRole(id, req.RoleID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)

	currentUser := infraMiddleware.GetUsername(c)
	h.events.SendNotification(currentUser, "角色分配", "为用户 #"+strconv.FormatUint(uint64(id), 10)+" 分配了新角色")
}

// @Summary      批量分配角色
// @Description  为一组用户批量分配角色
// @Tags         用户-角色
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.BatchAssignRoleRequest  true  "批量分配信息"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /users/batch/role [post]
func (h *UserHandler) BatchAssignRole(c *gin.Context) {
	var req dto.BatchAssignRoleRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.userOrchestration.BatchAssignRole(&orchestration.BatchAssignRoleRequest{
		UserIDs: req.UserIDs,
		RoleID:  req.RoleID,
	}); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)

	count := len(req.UserIDs)
	currentUser := infraMiddleware.GetUsername(c)
	h.events.SendNotification(currentUser, "批量分配角色", "为 "+strconv.Itoa(count)+" 名用户分配了角色")
}

// GetUsersByRole godoc
// @Summary      按角色查询用户
// @Description  根据角色 ID 查询所有该角色的用户
// @Tags         用户-角色
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        roleId  path  uint  true  "角色 ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /users/by-role/{roleId} [get]
func (h *UserHandler) GetUsersByRole(c *gin.Context) {
	roleIDStr := c.Param("roleId")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid role id")
		return
	}

	users, err := h.userOrchestration.GetUsersByRole(uint(roleID))
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	userList := make([]dto.UserByRoleItem, len(users))
	for i, u := range users {
		userList[i] = dto.UserByRoleItem{
			ID:       u.ID,
			Username: u.Username,
			Nickname: u.Nickname,
			Email:    u.Email,
			Status:   u.Status,
		}
	}

	utils.Success(c, userList)
}

// GetCurrentMenus godoc
// @Summary      获取当前用户菜单树
// @Description  返回当前用户角色拥有的菜单权限树形结构
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Router       /auth/menus [get]
func (h *UserHandler) GetCurrentMenus(c *gin.Context) {
	userID := infraMiddleware.GetUserID(c)
	menus, err := h.userOrchestration.GetCurrentUserMenus(userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, menus)
}

// GetCurrentPermissions godoc
// @Summary      获取当前用户权限标识
// @Description  返回当前用户角色拥有的所有权限标识列表（如 ["user:create", "role:delete"]）
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response{data=[]string}
// @Router       /auth/permissions [get]
func (h *UserHandler) GetCurrentPermissions(c *gin.Context) {
	userID := infraMiddleware.GetUserID(c)
	codes, err := h.userOrchestration.GetCurrentUserPermissions(userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.Success(c, codes)
}

// GetProfile godoc
// @Summary      获取个人信息
// @Description  获取当前登录用户的个人信息
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  utils.Response
// @Router       /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := infraMiddleware.GetUserID(c)
	user, err := h.userOrchestration.GetUser(userID)
	if err != nil {
		utils.NotFound(c, "user not found")
		return
	}
	resp := dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		Status:    user.Status,
		RoleID:    user.RoleID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	if user.Role != nil {
		resp.Role = &dto.RoleSimpleResponse{
			ID:   user.Role.ID,
			Name: user.Role.Name,
			Code: user.Role.Code,
		}
	}
	utils.Success(c, resp)
}

// UpdateProfile godoc
// @Summary      编辑个人信息
// @Description  当前用户编辑自己的昵称、邮箱、手机号
// @Tags         个人中心
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  object{nickname=string,email=string,phone=string}  true  "个人信息"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req dto.UpdateProfileRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	userID := infraMiddleware.GetUserID(c)
	if err := h.userOrchestration.UpdateProfile(userID, &orchestration.UpdateProfileRequest{
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
	}); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

// UpdateAvatar godoc
// @Summary      上传头像
// @Description  上传当前用户的头像图片
// @Tags         个人中心
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        avatar  formData  file  true  "头像文件（支持 jpg/png/gif，最大 5MB）"
// @Success      200  {object}  utils.Response{data=object{url=string}}
// @Failure      400  {object}  utils.Response
// @Router       /users/avatar [post]
func (h *UserHandler) UpdateAvatar(c *gin.Context) {
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		utils.BadRequest(c, "avatar file is required")
		return
	}
	defer file.Close()

	uploadCfg := config.GetConfig().Upload
	ext := filepath.Ext(header.Filename)
	allowed := false
	for _, e := range uploadCfg.AllowedExts {
		if e == ext {
			allowed = true
			break
		}
	}
	if !allowed {
		utils.BadRequest(c, "file type not allowed, only jpg/png/gif")
		return
	}

	data := make([]byte, 512)
	n, _ := file.Read(data)
	mimeType := http.DetectContentType(data[:n])
	if !allowedMimeTypes[mimeType] {
		utils.BadRequest(c, "invalid file content, only jpg/png/gif images are allowed")
		return
	}
	file.Seek(0, 0)

	if header.Size > int64(uploadCfg.MaxSize)*1024*1024 {
		utils.BadRequest(c, fmt.Sprintf("file too large, max %dMB", uploadCfg.MaxSize))
		return
	}

	avatarDir := filepath.Join(uploadCfg.Dir, "avatars")
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		utils.InternalError(c, "failed to create upload dir")
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(header.Filename))
	savePath := filepath.Join(avatarDir, filename)

	dst, err := os.Create(savePath)
	if err != nil {
		utils.InternalError(c, "failed to save file")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		utils.InternalError(c, "failed to save file")
		return
	}

	avatarURL := "/uploads/avatars/" + filename
	userID := infraMiddleware.GetUserID(c)

	utils.Success(c, dto.AvatarResponse{URL: avatarURL})

	go func() {
		if err := h.userOrchestration.UpdateAvatar(userID, avatarURL); err != nil {
			logger.Error("failed to update avatar in background",
				logger.UintField("user_id", userID),
				logger.ErrorField(err),
			)
		}
	}()
}

// UpdateStatus godoc
// @Summary      启用/禁用用户
// @Description  管理员启用或禁用指定用户账户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path  uint                    true  "用户 ID"
// @Param        body  body  object{status=int}  true  "状态（1=启用，2=禁用）"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Router       /users/{id}/status [put]
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	id, ok := utils.ParseID(c)
	if !ok {
		return
	}

	var req dto.UpdateStatusRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid status, must be 1(active) or 2(inactive)")
		return
	}

	if err := h.userOrchestration.UpdateStatus(id, req.Status); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)

	statusText := "启用"
	if req.Status == int(entity.UserStatusInactive) {
		statusText = "禁用"
	}
	currentUser := infraMiddleware.GetUsername(c)
	h.events.SendNotification(currentUser, "用户状态变更", "已将用户 #"+strconv.FormatUint(uint64(id), 10)+" "+statusText)
}

// ImportUsers godoc
// @Summary      批量导入用户
// @Description  通过 Excel 文件批量导入用户
// @Tags         用户管理
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        file  formData  file  true  "Excel 文件（表头：用户名,密码,昵称,邮箱,手机号,角色ID）"
// @Success      200  {object}  utils.Response{data=object{count=int}}
// @Failure      400  {object}  utils.Response
// @Router       /users/import [post]
func (h *UserHandler) ImportUsers(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "file is required")
		return
	}

	file, err := fh.Open()
	if err != nil {
		utils.BadRequest(c, "failed to open file")
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		utils.BadRequest(c, "invalid excel file")
		return
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil || len(rows) < 2 {
		utils.BadRequest(c, "empty file or missing header")
		return
	}

	var users []*orchestration.CreateUserRequest
	for i := 1; i < len(rows); i++ {
		cols := rows[i]
		if len(cols) < 2 {
			continue
		}

		username := strings.TrimSpace(cols[0])
		password := strings.TrimSpace(cols[1])
		if username == "" {
			utils.BadRequest(c, fmt.Sprintf("row %d: username is empty", i+1))
			return
		}
		if len(username) < 3 || len(username) > 64 {
			utils.BadRequest(c, fmt.Sprintf("row %d: username length must be 3-64 characters", i+1))
			return
		}
		if len(password) < 6 {
			utils.BadRequest(c, fmt.Sprintf("row %d: password length must be at least 6 characters", i+1))
			return
		}

		nickname := ""
		if len(cols) > 2 {
			nickname = cols[2]
		}
		email := ""
		if len(cols) > 3 {
			email = cols[3]
		}
		phone := ""
		if len(cols) > 4 {
			phone = cols[4]
		}
		roleID := uint(0)
		if len(cols) > 5 {
			parsed, err := strconv.ParseUint(cols[5], 10, 64)
			if err != nil || parsed == 0 {
				utils.BadRequest(c, fmt.Sprintf("row %d: invalid role_id", i+1))
				return
			}
			roleID = uint(parsed)
		} else {
			utils.BadRequest(c, fmt.Sprintf("row %d: role_id is required", i+1))
			return
		}

		users = append(users, &orchestration.CreateUserRequest{
			Username: username,
			Password: password,
			Nickname: nickname,
			Email:    email,
			Phone:    phone,
			RoleID:   roleID,
		})
	}

	if len(users) == 0 {
		utils.BadRequest(c, "no valid user data found")
		return
	}

	if err := h.userOrchestration.ImportUsers(users); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, dto.ImportResult{Count: len(users)})

	h.events.BroadcastStats()
	count := len(users)
	currentUser := infraMiddleware.GetUsername(c)
	h.events.SendNotification(currentUser, "批量导入", "成功导入 "+strconv.Itoa(count)+" 名用户")
}

// ExportUsers godoc
// @Summary      批量导出用户
// @Description  导出用户列表为 Excel 文件（分页流式写入）
// @Tags         用户管理
// @Accept       json
// @Produce      application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security     BearerAuth
// @Success      200  {file}  binary
// @Router       /users/export [get]
func (h *UserHandler) ExportUsers(c *gin.Context) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := f.GetSheetName(0)
	headers := []string{"用户名", "昵称", "邮箱", "手机号", "状态", "角色ID"}
	for i, header := range headers {
		col := string(rune('A' + i))
		f.SetCellValue(sheet, fmt.Sprintf("%s1", col), header)
	}

	type pageResult struct {
		users []entity.User
		page  int
	}
	fetchCh := make(chan pageResult, 4)
	errCh := make(chan error, 1)

	go func() {
		page := 1
		pageSize := 500
		defer close(fetchCh)
		for {
			users, total, err := h.userOrchestration.ListUsers(page, pageSize, nil)
			if err != nil {
				errCh <- err
				return
			}
			fetchCh <- pageResult{users: users, page: page}
			if int64(page*pageSize) >= total {
				return
			}
			page++
		}
	}()

	row := 2
	for res := range fetchCh {
		for _, u := range res.users {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), u.Username)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), u.Nickname)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), u.Email)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), u.Phone)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), u.Status)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), u.RoleID)
			row++
		}
	}

	select {
	case err := <-errCh:
		utils.InternalError(c, err.Error())
		return
	default:
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=users.xlsx")
	if err := f.Write(c.Writer); err != nil {
		utils.InternalError(c, "failed to export")
	}
}
