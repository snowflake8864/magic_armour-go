package port

type AddPolicy struct {
	Protocol    string `form:"protocol" json:"protocol" binding:"required,oneof=tcp"`
	Type        string `form:"type" json:"type" binding:"required,oneof=ipv4 ipv6"`
	SourcePort  string `form:"sourcePort" json:"sourcePort" binding:"required"`
	ForwardIp   string `form:"forwardIp" json:"forwardIp" binding:"omitempty,ip"`
	ForwardType *int   `form:"forwardType" json:"forwardType" binding:"required,oneof=0 1"`
	ForwardPort int    `form:"forwardPort" json:"forwardPort" binding:"omitempty,min=1,max=65535"`
	AlarmLevel  int    `form:"alarmLevel" json:"alarmLevel" binding:"required,min=1,max=5"`
}

type Id struct {
	Id string `form:"id" json:"id" binding:"required,min=1"`
}

type DeptId struct {
	DeptId int `form:"deptId" json:"deptId"`
}

type SearchPort struct {
	Port int `form:"port" json:"port" binding:"omitempty,min=1,max=65535"`
}
