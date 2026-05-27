package dto

type CreatePermissionRequest struct {
	Name     string `json:"name" binding:"required,max=64"`
	Code     string `json:"code" binding:"required,max=64"`
	Type     int    `json:"type" binding:"oneof=1 2"`
	ParentID uint   `json:"parent_id"`
	Path     string `json:"path" binding:"max=256"`
	Method   string `json:"method" binding:"max=32"`
	Sort     int    `json:"sort"`
}

type UpdatePermissionRequest struct {
	Name     string `json:"name" binding:"required,max=64"`
	Code     string `json:"code" binding:"required,max=64"`
	Type     int    `json:"type" binding:"oneof=1 2"`
	ParentID uint   `json:"parent_id"`
	Path     string `json:"path" binding:"max=256"`
	Method   string `json:"method" binding:"max=32"`
	Sort     int    `json:"sort"`
}
