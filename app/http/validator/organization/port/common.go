package port

type AddPolicy struct {
	Protocol    string `form:"protocol" json:"protocol" binding:"required,oneof=tcp"`
	Type        string `form:"type" json:"type" binding:"required,oneof=ipv4 ipv6"`
	SourcePort  string `form:"sourcePort" json:"sourcePort" binding:"required"`
	AlarmLevel  int    `form:"alarmLevel" json:"alarmLevel" binding:"required,min=1,max=5"`
	ForwardIp   string `form:"forwardIp" json:"forwardIp" binding:"omitempty,ip"`
	ForwardPort int    `form:"forwardPort" json:"forwardPort" binding:"omitempty,min=1,max=65535"`
}

type Id struct {
	Id int64 `form:"id" json:"id" binding:"required,min=1"`
}

type DeptId struct {
	DeptId int `form:"deptId" json:"deptId"`
}

type Suffix struct {
	Suffix string `form:"suffix" json:"suffix" binding:"required,max=64"`
}
