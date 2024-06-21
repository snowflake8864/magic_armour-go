package alarm

type Id struct {
	Id int64 `form:"id" json:"id" binding:"required,min=1"`
}

type Remark struct {
	Remark string `form:"remark" json:"remark" binding:"required,max=128"`
}

type Hash struct {
	Hash string `form:"hash" json:"hash" binding:"omitempty,max=128"`
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
