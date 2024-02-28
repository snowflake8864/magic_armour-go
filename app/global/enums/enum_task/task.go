package enum_task

// 任务

const (
	StatusStart  int = 0
	StatusFinish int = 1

	TypeProcessRefresh            int = 0  //任务类型: 刷新进程
	TypeAgentUpgrade              int = 1  //任务类型:软件升级
	TypeDirStructureUpload        int = 2  //任务类型:上传目录结构
	TypeProcessWhiteDownload      int = 3  //任务类型:下载进程白名单库
	TypeTamperPolicyDownload      int = 4  //任务类型:下载防篡改策略
	TypelistUploadConfig          int = 5  //任务类型:上传当前配置信息 （没有使用）
	TypeAgentConfigDownload       int = 6  //任务类型:下载终端配置信息
	TypeProcessBlackDownload      int = 7  //任务类型:下载进程黑名单库
	TypePortWorkingUpload         int = 9  //任务类型:上传正常业务端口
	TypePortVirtualPolicyDownload int = 10 //任务类型:下载端口虚开策略
	TypeDownloadPluggingConfig    int = 11 //任务类型:下载封堵策略（不在使用）
	TypeUploadPluggingPort        int = 12 //任务类型: 上传自动封堵列表（不在使用）
	TypePluggingDownload          int = 13 //任务类型:下载封堵列表
	TypePlugIpWhiteDownload       int = 14 //任务类型:下载封堵ip白名单（不在使用）
	TypePlugIpBlackDownload       int = 15 //任务类型:下载封堵ip黑名单
	TypeUsbUpload                 int = 16 //任务类型: 上传外设列表
	TypeBlackmailPolicyDownload   int = 19 //任务类型:  下载勒索防护策略
	TypeProcessCollectDll         int = 21 //任务类型: 收集指定进程并相关模块
	TypeDllRefresh                int = 22 //任务类型: 刷新模块
	TypeDllWhiteDownload          int = 23 //任务类型: 下载模块白名单
	TypeDllBlackDownload          int = 24 //任务类型: 下载模块黑名单
	TypeAgentUninstall            int = 25 //任务类型:软件卸载
	TypeUsbWhiteDownload          int = 26 //任务类型:外设信任名单
	TypeUsbBlackDownload          int = 27 //任务类型:外设禁止名单
	TypeProcessSampleUpload       int = 28 //任务类型:进程样本提取
	TypeSyslogOpenSwitch          int = 29 //任务类型: 开启syslog（不在使用）
	TypeSyslogCloseSwitch         int = 30 //任务类型: 关闭 syslog（不在使用）
	TypeProcessSuffixIsExe        int = 31 //任务类型:获取所有程序
	TypeDllSuffixIsDll            int = 32 //任务类型:获取所有模块
	TypeTrustDir                  int = 33 //任务类型:信任目录

)

var EnumMap = map[string]map[int]string{
	"status": {
		StatusStart:  "进行中",
		StatusFinish: "已完成",
	},
	"type": {
		TypeAgentConfigDownload:       "下发终端配置",
		TypeTrustDir:                  "下发信任目录",
		TypeProcessRefresh:            "刷新运行中进程",
		TypeProcessSuffixIsExe:        "获取所有程序",
		TypeProcessWhiteDownload:      "下发进程白名单",
		TypeProcessBlackDownload:      "下发进程黑名单",
		TypeProcessCollectDll:         "收集指定进程下模块",
		TypeProcessSampleUpload:       "提取样本",
		TypeDllRefresh:                "刷新运行中模块",
		TypeDllSuffixIsDll:            "获取所有模块",
		TypeDllWhiteDownload:          "下发模块白名单",
		TypeDllBlackDownload:          "下发模块黑名单",
		TypeDirStructureUpload:        "获取目录结构",
		TypeTamperPolicyDownload:      "下发防篡改策略",
		TypeBlackmailPolicyDownload:   "下发勒索防护策略",
		TypePortWorkingUpload:         "获取系统正常业务端口",
		TypePortVirtualPolicyDownload: "下发端口虚开策略",
		TypePluggingDownload:          "下发动态封堵策略",
		TypePlugIpBlackDownload:       "下发封堵ip黑名单",
		TypeUsbUpload:                 "获取外设列表",
		TypeUsbWhiteDownload:          "下发外设信任名单",
		TypeUsbBlackDownload:          "下发外设禁止名单",
		TypeAgentUpgrade:              "升级终端",
		TypeAgentUninstall:            "卸载终端",
	},
}

var EnumSlice = map[string][]int{
	//windows xp 不支持的任务
	"win_xp_not_support": {
		TypePortWorkingUpload,
		TypePortVirtualPolicyDownload,
		TypePluggingDownload,
		TypePlugIpBlackDownload,
	},
	//linux 不支持的任务
	"linux_not_support": {
		TypeDllRefresh,
		TypeDllWhiteDownload,
		TypeDllBlackDownload,
		TypeDllSuffixIsDll,
		TypeProcessCollectDll,
	},
	//在保护模式(配置中有开启保护模式)之前下发的任务
	"before_project_mode": {
		TypeProcessWhiteDownload,
		TypeProcessBlackDownload,
		TypeDllWhiteDownload,
		TypeDllBlackDownload,
		TypeTamperPolicyDownload,
		TypeBlackmailPolicyDownload,
		TypeUsbWhiteDownload,
		TypeUsbBlackDownload,
	},
}
