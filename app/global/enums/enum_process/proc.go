package enum_process

// 进程或模块

const (
	SourceProcess int = 1 //数据来源:进程
	SourceDll     int = 2 //数据来源:模块

	StatusUnset   int = 1 //状态： 未设置(默认)
	StatusUnknown int = 2 //状态： 未知
	StatusWhite   int = 4 //状态： 白名单
	StatusBlack   int = 8 //状态： 黑名单

	RelyStatusUnset int = 0 //(self_status,threaten_status,global_status)状态： 未设置
	RelyStatusWhite int = 1 //(self_status,threaten_status,global_status)状态： 白名单
	RelyStatusBlack int = 2 //(self_status,threaten_status,global_status)状态： 黑名单

	SearchStatusUnset         int = 1   //状态： 未设置
	SearchStatusUnknown       int = 2   //状态： 未知
	SearchStatusWhite         int = 4   //状态： 白名单
	SearchStatusBlack         int = 8   //状态： 黑名单
	SearchStatusGlobalWhite   int = 16  //状态： 全局白名单
	SearchStatusGlobalBlack   int = 32  //状态： 全局黑名单
	SearchStatusThreatenWhite int = 64  //状态： 威胁情报白名单
	SearchStatusThreatenBlack int = 128 //状态： 威胁情报黑名单

	ProcGetDllStatusDefault int = 1 //进程获取指定模块状态： 默认
	ProcGetDllStatusBegin   int = 2 // 进程获取指定模块状态： 开始(等待终端获取任务)
	ProcGetDllStatusRun     int = 3 // 进程获取指定模块状态： 获取中(等待终端上传数据)
	ProcGetDllStatusFinish  int = 4 //进程获取指定模块状态：结束

	SampleStatusDefault  int = 1 //样本状态： 默认
	SampleStatusGetting  int = 2 //样本状态：提取中
	SampleStatusFinished int = 3 //样本状态：已完成
)

var EnumMapInt = map[string]map[int]string{
	"search_status": {
		SearchStatusUnset:         "未设置",
		SearchStatusUnknown:       "未知",
		SearchStatusWhite:         "白名单",
		SearchStatusBlack:         "黑名单",
		SearchStatusGlobalWhite:   "全局白名单",
		SearchStatusGlobalBlack:   "全局黑名单",
		SearchStatusThreatenWhite: "威胁情报白名单",
		SearchStatusThreatenBlack: "威胁情报黑名单",
	},
}

//
