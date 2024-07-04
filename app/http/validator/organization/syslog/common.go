package syslog

type Ids struct {
	Ids string `form:"ids" json:"ids" binding:"required,min=6"`
}

type ProdHash struct {
	ProdHash string `form:"prodHash" json:"prodHash" binding:"omitempty,max=128"`
}

type Process struct {
	ProdId      int    `form:"prodId" json:"prodId" binding:"omitempty,max=65535"`
	ProdDir     string `form:"prodDir" json:"prodDir" binding:"omitempty,max=65535"`
	ProcParam   string `form:"procParam" json:"procParam" binding:"omitempty,min=1,max=256"`
	ParentHash  string `form:"parentHash" json:"parentHash" binding:"omitempty,max=128"`
	ParentId    int    `form:"parentId" json:"parentId" binding:"omitempty,min=1,max=65535"`
	ParentDir   string `form:"parentDir" json:"parentDir" binding:"omitempty,min=1,max=65535"`
	ParentParam string `form:"parentParam" json:"parentParam" binding:"omitempty,min=1,max=256"`
}

type DomainName struct {
	DomainName string `form:"domainName" json:"domainName" binding:"omitempty,max=128"`
}

type SourceIp struct {
	SourceIp string `form:"sourceIp" json:"sourceIp" binding:"omitempty,ip"`
}

type SourcePort struct {
	SourcePort int `form:"sourcePort" json:"sourcePort" binding:"omitempty,max=65535"`
}

type TargetIp struct {
	TargetIp string `form:"targetIp" json:"targetIp" binding:"omitempty,ip"`
}
type TargetPort struct {
	TargetPort int `form:"targetPort" json:"targetPort" binding:"omitempty,max=65535"`
}

type LoginStatus struct {
	LoginStatus *int `form:"loginStatus" json:"loginStatus" binding:"required,oneof=-100 0 1"`
}

type SearchTime struct {
	StartTime string `form:"startTime" json:"startTime" binding:"required,datetime=2006-01-02 15:04:05"`
	EndTime   string `form:"endTime" json:"endTime" binding:"required,datetime=2006-01-02 15:04:05"`
}
