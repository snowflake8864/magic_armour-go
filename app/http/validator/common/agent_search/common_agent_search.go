package agent_search

// Conditions 终端条件
type Conditions struct {
	Page        int    `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize    int    `form:"pageSize" json:"pageSize" binding:"omitempty,min=1"`
	Field       string `form:"field" json:"field" binding:"omitempty,min=1,max=32"`
	Order       string `form:"order" json:"order" binding:"omitempty,oneof=asc desc"`
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
