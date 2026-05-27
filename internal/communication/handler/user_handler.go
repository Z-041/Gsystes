package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/communication/dto"
	"github.com/gsystes/backend/internal/infrastructure/auth"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
	orchestration "github.com/gsystes/backend/internal/orchestration/service"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

type UserHandler struct {
	userOrchestration *orchestration.UserOrchestration
}

func NewUserHandler(userOrchestration *orchestration.UserOrchestration) *UserHandler {
	return &UserHandler{userOrchestration: userOrchestration}
}

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
	}, auth.GenerateToken)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.Success(c, gin.H{
		"token": resp.Token,
		"user": gin.H{
			"id":       resp.User.ID,
			"username": resp.User.Username,
			"nickname": resp.User.Nickname,
			"avatar":   resp.User.Avatar,
			"role_id":  resp.User.RoleID,
		},
	})
}

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

	utils.Success(c, gin.H{
		"id": user.ID,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	var req dto.UpdateUserRequest
	valErrors := utils.BindAndValidate(c, &req)
	if valErrors != nil {
		utils.BadRequest(c, "invalid request parameters")
		return
	}

	if err := h.userOrchestration.UpdateUser(&orchestration.UpdateUserRequest{
		ID:       uint(id),
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

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	if err := h.userOrchestration.DeleteUser(uint(id)); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil)
}

func (h *UserHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.BadRequest(c, "invalid user id")
		return
	}

	user, err := h.userOrchestration.GetUser(uint(id))
	if err != nil {
		utils.NotFound(c, "user not found")
		return
	}

	utils.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"email":      user.Email,
		"phone":      user.Phone,
		"avatar":     user.Avatar,
		"status":     user.Status,
		"role_id":    user.RoleID,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func (h *UserHandler) List(c *gin.Context) {
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

	conditions := make(map[string]interface{})
	if username := c.Query("username"); username != "" {
		conditions["username LIKE ?"] = "%" + username + "%"
	}
	if status := c.Query("status"); status != "" {
		conditions["status = ?"] = status
	}

	users, total, err := h.userOrchestration.ListUsers(page, pageSize, conditions)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	userList := make([]gin.H, len(users))
	for i, u := range users {
		userList[i] = gin.H{
			"id":         u.ID,
			"username":   u.Username,
			"nickname":   u.Nickname,
			"email":      u.Email,
			"phone":      u.Phone,
			"avatar":     u.Avatar,
			"status":     u.Status,
			"role_id":    u.RoleID,
			"created_at": u.CreatedAt,
		}
	}

	utils.PageSuccess(c, userList, total, page, pageSize)
}

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