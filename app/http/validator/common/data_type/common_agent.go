package data_type

// AgentSearch 终端查询条件
type AgentSearch struct {
	AgentId     string `form:"agentId" json:"agentId" binding:"omitempty,min=8,max=32"`
	DeptId      int    `form:"deptId" json:"deptId"`
	HostName    string `form:"hostName" json:"hostName" binding:"omitempty,min=1"`
	AgentIp     string `form:"agentIp" json:"agentIp" binding:"omitempty,ip"`
	NetIp       string `form:"netIp" json:"netIp" binding:"omitempty,ip"`
	Alias       string `form:"alias" json:"alias" binding:"omitempty,min=1"`
	ContactName string `form:"contactName" json:"contactName" binding:"omitempty,min=1"`
	IsOnline    *int   `form:"isOnline" json:"isOnline" binding:"omitempty,oneof=-100 0 1"`
	IsSafe      *int   `form:"isSafe" json:"isSafe" binding:"omitempty,oneof=-100 0 1"`
	GrantStatus *int   `form:"grantStatus" json:"grantStatus" binding:"omitempty,oneof=-100 0 1"`
}

// AgentSingle 终端单个
type AgentSingle struct {
	AgentId string `form:"agentId" json:"agentId" binding:"required,min=8,max=32"`
}

// AgentBatch 终端批量
type AgentBatch struct {
	AgentIds []string `form:"agentIds" json:"agentIds" binding:"required,gt=0,dive,min=8,max=32"`
}

// AgentSingleOmitempty 终端单个
type AgentSingleOmitempty struct {
	AgentId string `form:"agentId" json:"agentId" binding:"omitempty,min=8,max=32"`
}
