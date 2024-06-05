package proc_global

type AddSingle struct {
	Id   string `form:"id" json:"id" binding:"required,min=8,max=32"`
	Dir  string `form:"dir" json:"dir" binding:"required,max=512"`
	Type int    `form:"type" json:"type" binding:"required,oneof=1 2"`
}

type AddBatch struct {
	List string `form:"list" json:"list" binding:"required,min=32"`
}

type PageSearch struct {
	Id   string `form:"id" json:"id" binding:"omitempty,min=8,max=32"`
	Dir  string `form:"dir" json:"dir" binding:"omitempty,max=512"`
	Type int    `form:"type" json:"type" binding:"required,oneof=1 2"`
}

type IdSingle struct {
	Id string `form:"id" json:"id" binding:"required,min=8,max=32"`
}
