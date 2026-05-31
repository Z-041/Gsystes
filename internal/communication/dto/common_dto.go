package dto

type PageParam struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

type IDParam struct {
	ID uint `uri:"id" binding:"required"`
}

type MenuTreeNode struct {
	ID       uint            `json:"id"`
	Name     string          `json:"name"`
	Code     string          `json:"code"`
	Path     string          `json:"path"`
	Sort     int             `json:"sort"`
	Children []*MenuTreeNode `json:"children"`
}

type IDResponse struct {
	ID uint `json:"id"`
}
