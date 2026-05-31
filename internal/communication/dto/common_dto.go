package dto

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
