package role

type Add struct {
	Name   string `form:"name" json:"name" binding:"required,max=24"`
	Remark string `form:"remark" json:"remark" binding:"omitempty,max=128"`
}
