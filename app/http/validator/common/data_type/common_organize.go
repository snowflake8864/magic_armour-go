package data_type

// OrganizeSearch 企业查询条件
type OrganizeSearch struct {
	OrgType   string `form:"orgType" json:"orgType" binding:"required"`
	OraName   string `form:"oraName" json:"oraName" binding:"omitempty,min=2,max=128"`
	UserName  string `form:"userName" json:"userName" binding:"omitempty,min=2,max=128"`
	UserPhone string `form:"userPhone" json:"userPhone" binding:"omitempty,min=8,max=11"`
}

// OrganizeSingle 企业单个
type OrganizeSingle struct {
	OrgId int `form:"orgId" json:"orgId" binding:"required"`
}
