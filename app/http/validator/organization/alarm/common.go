package alarm

type Id struct {
	Id string `form:"id" json:"id" binding:"required,min=1"`
}

type Remark struct {
	Remark string `form:"remark" json:"remark" binding:"required,max=128"`
}

type Hash struct {
	Hash string `form:"hash" json:"hash" binding:"required,max=128"`
}

type SearchHash struct {
	Hash string `form:"hash" json:"hash" binding:"omitempty,max=128"`
}

type SearchFileDir struct {
	FileDir string `form:"fileDir" json:"fileDir" binding:"omitempty,max=128"`
}

type AttackIp struct {
	AttackIp string `form:"attackIp" json:"attackIp" binding:"required,ip"`
}

type SearchAttackIp struct {
	AttackIp string `form:"attackIp" json:"attackIp" binding:"omitempty,ip"`
}

type PlugIp struct {
	PlugIp string `form:"plugIp" json:"plugIp" binding:"required,ip"`
}

type SearchPlugIp struct {
	PlugIp string `form:"plugIp" json:"plugIp" binding:"omitempty,ip"`
}

type SourceIp struct {
	SourceIp string `form:"sourceIp" json:"sourceIp" binding:"required,ip"`
}

type SearchSourceIp struct {
	SourceIp string `form:"sourceIp" json:"sourceIp" binding:"omitempty,ip"`
}

type TargetIp struct {
	TargetIp string `form:"targetIp" json:"targetIp" binding:"required,ip"`
}

type SearchTargetIp struct {
	TargetIp string `form:"targetIp" json:"targetIp" binding:"omitempty,ip"`
}

type HandleStatus struct {
	HandleStatus int `form:"handleStatus" json:"handleStatus" binding:"required,oneof=-100 1 2"`
}

type SearchType struct {
	AlarmType int `form:"alarmType" json:"alarmType" binding:"required"`
}

type SearchTime struct {
	StartTime string `form:"startTime" json:"startTime" binding:"required,datetime=2006-01-02 15:04:05"`
	EndTime   string `form:"endTime" json:"endTime" binding:"required,datetime=2006-01-02 15:04:05"`
}
