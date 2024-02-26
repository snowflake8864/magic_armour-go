package enum_agent

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

)
