package usb

type StatusSet struct {
	Status int `form:"status" json:"status" binding:"required,oneof=1 2"`
}

type Id struct {
	Id string `form:"id" json:"id" binding:"required,min=16"`
}

type AliasSet struct {
	Alias  string `form:"alias" json:"alias" binding:"required,max=128"`
	Remark string `form:"remark" json:"remark" binding:"omitempty,max=128"`
}
