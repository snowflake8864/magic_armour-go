package common_data

// ProcUpload0 进程上传
type ProcUpload0 struct {
	ProcUpload string `form:"proclist" json:"proclist" binding:"required"`
}

// ProcUpload1 进程上传
type ProcUpload1 struct {
	ProcUpload string `form:"list" json:"list" binding:"required"`
}

// AlarmUpload 告警上传
type AlarmUpload struct {
	Alert string `form:"alert" json:"alert" binding:"required"`
}

// SyslogUpload 通讯日志上传
type SyslogUpload struct {
	List string `form:"list" json:"list" binding:"required"`
}
