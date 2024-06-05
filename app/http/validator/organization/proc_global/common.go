package proc_global

type AddSingle struct {
	Hash string `form:"hash" json:"hash" binding:"required,min=8,max=32"`
	Dir  string `form:"dir" json:"dir" binding:"required,max=512"`
	Type int    `form:"type" json:"type" binding:"required,oneof=1 2"`
}
