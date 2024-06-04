package proc_sample

type Add struct {
	Hash      string `form:"hash" json:"hash" binding:"required,min=8,max=32"`
	Dir       string `form:"dir" json:"dir" binding:"required,max=512"`
	Introduce string `form:"introduce" json:"introduce" binding:"omitempty,max=256"`
	Copyright string `form:"copyright" json:"copyright" binding:"omitempty,max=256"`
	Signature string `form:"signature" json:"signature" binding:"omitempty,max=512"`
	Attribute string `form:"attribute" json:"attribute" binding:"omitempty,max=1024"`
}
