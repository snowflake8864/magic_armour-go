package data_type

type Page struct {
	Page     int    `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"pageSize" json:"pageSize" binding:"omitempty,min=1,max=2000"`
	Field    string `form:"field" json:"field" binding:"omitempty,min=1,max=32"`
	Order    string `form:"order" json:"order" binding:"omitempty,oneof=asc desc"`
}
