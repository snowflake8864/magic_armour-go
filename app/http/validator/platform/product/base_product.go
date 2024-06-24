package product

// PAdd 产品
type PAdd struct {
	GrantType   int     `form:"grant_type" json:"grant_type" binding:"required,oneof=1 2"`
	Title       string  `form:"title" json:"title" binding:"required,min=4,max=128"`
	GrantModule string  `form:"grant_module" json:"grant_module" binding:"required,min=10"`
	Price       float64 `form:"price"   json:"price" price:"required,max=100000"`
	Remark      string  `form:"remark" json:"remark" binding:"omitempty,max=200"`
}
