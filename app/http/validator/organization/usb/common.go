package usb

type StatusSet struct {
	Status int `form:"status" json:"status" binding:"required,oneof=2 3"`
}

type Id struct {
	Id string `form:"id" json:"id" binding:"required,min=6,max=32"`
}

type AliasSet struct {
	Alias  string `form:"alias" json:"alias" binding:"required,max=128"`
	Remark string `form:"remark" json:"remark" binding:"omitempty,max=128"`
}

type AddGlobal struct {
	Name  string `form:"name" json:"name" binding:"required,max=128"`
	Type  string `form:"type" json:"type" binding:"omitempty,max=256"`
	Intro string `form:"intro" json:"intro" binding:"omitempty,max=256"`
}
