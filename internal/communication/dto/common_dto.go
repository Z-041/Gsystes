package dto

type PageParam struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=1,max=100"`
}

type IDParam struct {
	ID uint `uri:"id" binding:"required"`
}
