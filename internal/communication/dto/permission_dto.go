package dto

import "time"

type CreatePermissionRequest struct {
	Name     string `json:"name" binding:"required,max=64"`
	Code     string `json:"code" binding:"required,max=64"`
	Type     int    `json:"type" binding:"required,oneof=1 2"`
	ParentID uint   `json:"parent_id"`
	Path     string `json:"path" binding:"max=256"`
	Method   string `json:"method" binding:"max=32"`
	Sort     int    `json:"sort"`
}

type UpdatePermissionRequest struct {
	Name     string `json:"name" binding:"required,max=64"`
	Code     string `json:"code" binding:"required,max=64"`
	Type     int    `json:"type" binding:"required,oneof=1 2"`
	ParentID uint   `json:"parent_id"`
	Path     string `json:"path" binding:"max=256"`
	Method   string `json:"method" binding:"max=32"`
	Sort     int    `json:"sort"`
}

type PermissionResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Type      int       `json:"type"`
	ParentID  uint      `json:"parent_id"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Sort      int       `json:"sort"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
