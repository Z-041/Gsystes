package dto

import "time"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Nickname string `json:"nickname" binding:"max=64"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	RoleID   uint   `json:"role_id" binding:"required"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname" binding:"max=64"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
	RoleID   uint   `json:"role_id"`
	Status   int    `json:"status" binding:"oneof=1 2"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
}

type BatchAssignRoleRequest struct {
	UserIDs []uint `json:"user_ids" binding:"required,min=1,dive,required"`
	RoleID  uint   `json:"role_id" binding:"required"`
}

type LoginResponse struct {
	Token string     `json:"token"`
	User  UserSimple `json:"user"`
}

type UserSimple struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	RoleID   uint   `json:"role_id"`
}

type UserResponse struct {
	ID        uint               `json:"id"`
	Username  string             `json:"username"`
	Nickname  string             `json:"nickname"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Avatar    string             `json:"avatar"`
	Status    int                `json:"status"`
	RoleID    uint               `json:"role_id"`
	Role      *RoleSimpleResponse `json:"role,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}

type UserListItem struct {
	ID        uint               `json:"id"`
	Username  string             `json:"username"`
	Nickname  string             `json:"nickname"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Avatar    string             `json:"avatar"`
	Status    int                `json:"status"`
	RoleID    uint               `json:"role_id"`
	Role      *RoleSimpleResponse `json:"role,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=64"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,max=20"`
}

type UpdateStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=1 2"`
}

type AssignRoleRequest struct {
	RoleID uint `json:"role_id" binding:"required"`
}

type UserByRoleItem struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}

type AvatarResponse struct {
	URL string `json:"url"`
}

type ImportResult struct {
	Count int `json:"count"`
}
