package data_type

// AgentSearch 终端查询条件
type AgentSearch struct {
	AgentId     string `form:"agentId" json:"agentId" binding:"omitempty,min=8,max=32"`
	DeptIds     []int  `form:"deptId" json:"deptId" binding:"omitempty,dive"`
	HostName    string `form:"hostName" json:"hostName" binding:"omitempty,min=1"`
	AgentIp     string `form:"agentIp" json:"agentIp" binding:"omitempty,ip"`
	NetIp       string `form:"netIp" json:"netIp" binding:"omitempty,ip"`
	Alias       string `form:"alias" json:"alias" binding:"omitempty,min=1"`
	ContactName string `form:"contactName" json:"contactName" binding:"omitempty,min=1"`
	IsOnline    int    `form:"isOnline" json:"isOnline"`
	IsSafe      int    `form:"isSafe" json:"isSafe"`
}

// AgentSingle 终端单个
type AgentSingle struct {
	AgentId string `form:"agentId" json:"agentId" binding:"required,min=8,max=32"`
}

// AgentBatch 终端批量
type AgentBatch struct {
	AgentIds []string `form:"agentIds" json:"agentIds" binding:"required,gt=0,dive,min=8,max=32"`
}

// AgentBasic 终端基础信息
type AgentBasic struct {
	DeptId      int64  `form:"deptId" json:"deptId" binding:"required,min=-1"`
	AgentType   int    `form:"agentType" json:"agentType" binding:"required,min=-1"`
	Alias       string `form:"alias" json:"alias" binding:"required,min=1"`
	ContactName string `form:"contactName" json:"contactName" binding:"required,min=1"`
	ContactTel  string `form:"contactTel" json:"contactTel" binding:"required,min=2"`
	Position    string `form:"position" json:"position" binding:"required,min=2"`
	Remark      string `form:"remark" json:"remark" binding:"required,min=2"`
}
