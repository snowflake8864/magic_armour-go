package enum_agent

import "goskeleton/app/global/enums"

// 终端

const (
	ModeMonitor int = 0 //模式：监控
	ModeProject int = 1 //模式：保护

	UpgradeStatusDefault     int = 0 //升级状态：默认
	upgradeStatusIssue       int = 1 //升级状态：任务进行中
	upgradeStatusDownloading int = 2 //升级状态：下载中
	upgradeStatusSuccess     int = 3 //升级状态:成功
	upgradeStatusFail        int = 4 //升级状态：失败

	HeartbeatCommonTime int = 60 //心跳:普通时间（秒）
	HeartbeatFastTime   int = 15 //心跳:快速时间（秒）

	HardwareTime int = 300 //上传硬件资源时间间隔（秒）

	LogProtoUdp int = 1 //日志协议:TDP
	LogProtoTcp int = 2 //日志协议:TCP

	OrderFieldDefault string = "id" //默认排序字段
)

var EnumIntMap = map[string]map[int]string{
	"agent_type": {
		1: "服务器",
		2: "业务系统",
		3: "办公终端",
		4: "笔记本",
		5: "台式机",
		6: "其它",
	},
	"is_online": {
		enums.IsYes: "在线",
		enums.IsNot: "离线",
	},
	"is_safe": {
		enums.IsYes: "正常",
		enums.IsNot: "异常",
	},
}

var WebMapInt = map[string][]int{
	"agent_type": {
		1,
		2,
		3,
		4,
		5,
		6,
	},
	"is_online": {
		enums.IsYes,
		enums.IsNot,
	},
	"is_safe": {
		enums.IsYes,
		enums.IsNot,
	},
}
