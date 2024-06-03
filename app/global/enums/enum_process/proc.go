package enum_process

// 进程或模块

const (
	SourceProcess int = 1 //数据来源:进程
	SourceDll     int = 2 //数据来源:模块

	ProcStatusUnset   int = 1 //状态： 未设置(默认)
	ProcStatusUnknown int = 2 //状态： 未知
	ProcStatusWhite   int = 4 //状态： 白名单
	ProcStatusBlack   int = 8 //状态： 黑名单

	RelyStatusUnset   int = 0 //(self_status,threaten_status,global_status)状态： 未设置
	RelyStatusWhite   int = 1 //(self_status,threaten_status,global_status)状态： 白名单
	RelyStatusBlack   int = 2 //(self_status,threaten_status,global_status)状态： 黑名单
	RelyStatusUnknown int = 3 //(self_status,threaten_status,global_status)状态： 未知

	AllStatusUnset         int = 1   //状态： 未设置
	AllStatusUnknown       int = 2   //状态： 未知
	AllStatusWhite         int = 4   //状态： 白名单
	AllStatusBlack         int = 8   //状态： 黑名单
	AllStatusGlobalWhite   int = 16  //状态： 全局白名单
	AllStatusGlobalBlack   int = 32  //状态： 全局黑名单
	AllStatusThreatenWhite int = 64  //状态： 威胁情报白名单
	AllStatusThreatenBlack int = 128 //状态： 威胁情报黑名单

	CollectDllStatusDefault int = 1 //进程获取指定模块状态： 默认
	CollectDllStatusBegin   int = 2 // 进程获取指定模块状态： 开始(等待终端获取任务)
	CollectDllStatusRun     int = 3 // 进程获取指定模块状态： 获取中(等待终端上传数据)
	CollectDllStatusFinish  int = 4 //进程获取指定模块状态：结束

	SampleStatusDefault  int = 1 //样本状态： 默认
	SampleStatusGetting  int = 2 //样本状态：提取中
	SampleStatusFinished int = 3 //样本状态：已完成
	SampleStatusFail     int = 4 //样本状态：已完成
)

var EnumIntMap = map[string]map[int]string{
	"all_status": {
		AllStatusUnset:         "未设置",
		AllStatusUnknown:       "未知",
		AllStatusWhite:         "白名单",
		AllStatusBlack:         "黑名单",
		AllStatusGlobalWhite:   "全局白名单",
		AllStatusGlobalBlack:   "全局黑名单",
		AllStatusThreatenWhite: "威胁情报白名单",
		AllStatusThreatenBlack: "威胁情报黑名单",
	},
	"sample_status": {
		SampleStatusDefault:  "未开始",
		SampleStatusGetting:  "提取中",
		SampleStatusFinished: "完成",
		SampleStatusFail:     "失败",
	},
}

var WebMapInt = map[string][]int{
	"proc_status": {
		ProcStatusUnset,
		ProcStatusUnknown,
		ProcStatusWhite,
		ProcStatusBlack,
	},
	"change_status": {
		AllStatusUnset,
		AllStatusUnknown,
		AllStatusWhite,
		AllStatusBlack,
		AllStatusGlobalWhite,
		AllStatusGlobalBlack,
	},
}

//
