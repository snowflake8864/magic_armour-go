package common_data

// OrganizeSearch 企业查询条件
type OrganizeSearch struct {
	OrgType   int    `form:"orgType" json:"orgType" binding:"required"`
	OraName   string `form:"oraName" json:"oraName" binding:"omitempty,min=2,max=128"`
	UserName  string `form:"userName" json:"userName" binding:"omitempty,min=2,max=128"`
	UserPhone string `form:"userPhone" json:"userPhone" binding:"omitempty,len=11"`
}

// OrganizeSingle 企业单个
type OrganizeSingle struct {
	OrgId int `form:"orgId" json:"orgId" binding:"required"`
}

// CommonReason 原因
type CommonReason struct {
	Reason string `form:"reason" json:"reason" binding:"required,max=200"`
}

// CommonSwitch 开关
type CommonSwitch struct {
	Switch *int `form:"switch" json:"switch" binding:"required,oneof=0 1"`
}
