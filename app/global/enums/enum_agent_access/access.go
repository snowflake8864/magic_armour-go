package enum_agent_access

// 终端

const (
	StatusSuccess int = 1
	StatusFail    int = 2
)

var EnumStrMap = map[string]map[string]string{
	"path": {
		"/v1/auth":      "登录",
		"/v1/getconf":   "下载配置信息",
		"/v1/gettask":   "获取任务(心跳)",
		"/v1/closetask": "关闭任务",

		"/v1/upload/suffix/exe":    "上传全部程序",
		"/v1/uploadproc":           "上传运行中进程",
		"/v1/autouploadprocess":    "主动上传进程",
		"/v1/getprocwl":            "下载进程白名单",
		"/v1/getprocbl":            "下载进程黑名单",
		"/v1/module/process":       "进程对应模块:下载指定进程",
		"/v1/upload/proc/module":   "进程对应模块:上传指定进程下的模块",
		"/v1/getdraw":              "查询未提取样本",
		"/v1/uploaddraw":           "上传样本",
		"/v1/upload/suffix/dll":    "上传全部模块",
		"/v1/pload/allproc/module": "上传运行中模块",
		"/v1/auto/upload/module":   "主动上传模块",
		"/v1/module/white/list":    "下载模块白名单",
		"/v1/module/black/list":    "下载模块黑名单",

		"/v1/getdirinfo":   "获取查询目录",
		"/v1/putdir":       "上传目录信息",
		"/v1/getdirpolicy": "下载防篡改策略",

		"/v1/getprotect": "下载勒索防护策略",

		"/v1/upserviceport":  "上传业务端口",
		"/v1/getvirtualport": "下载虚开端口",

		"/v1/getplugging":    "下载动态封堵策略",
		"/v1/getipblacklist": "下载封堵ip黑名单",

		"/v1/getOutreachDetect": "下载外联主动探测列表",

		"/v1/addperipherals":      "上传外设列表",
		"/v1/getwhiteperipherals": "获取外设信任（白名单）列表",
		"/v1/getblackperipherals": "获取外设禁止（黑名单）列表",

		"/v1/alertupload": "上传告警(除端口端口外)",
		"/v1/upopenport":  "上传端口虚开告警",

		"/v1/putsyslog": "上传通信日志",

		"/v1/puthardwareinfo": "上传终端资源信息",
		"/v1/uploadIp":        "上传终端IP信息",
		"/v1/download":        "升级",
		"/v1/uninstall":       "卸载",
		"/v1/error/upload":    "上传错误文件",
		"/v1/gettrustdir":     "获取信任目录",

		"/v1/getPwJump":      "获取口令跳变",
		"/v1/putPwJump":      "上报口令跳变结果",
		"/v1/getIpJump":      "获取IP跳变",
		"/v1/getBackups":     "获取备份策略",
		"/v1/uploadBackup":   "上传备份结果",
		"/v1/getRollbacks":   "获取回滚策略",
		"/v1/uploadRollback": "上传回滚结果",
	},
	"path_en": {
		"/v1/auth":      "Login",
		"/v1/getconf":   "Download Config",
		"/v1/gettask":   "Get Task (Heartbeat)",
		"/v1/closetask": "Close Task",

		"/v1/upload/suffix/exe":    "Upload All Programs",
		"/v1/uploadproc":           "Upload Running Processes",
		"/v1/autouploadprocess":    "Auto Upload Processes",
		"/v1/getprocwl":            "Download Process Whitelist",
		"/v1/getprocbl":            "Download Process Blacklist",
		"/v1/module/process":       "Process Modules: Download Specified Process",
		"/v1/upload/proc/module":   "Process Modules: Upload Modules for Specified Process",
		"/v1/getdraw":              "Query Un-extracted Samples",
		"/v1/uploaddraw":           "Upload Sample",
		"/v1/upload/suffix/dll":    "Upload All Dlls",
		"/v1/pload/allproc/module": "Upload Running Dlls",
		"/v1/auto/upload/module":   "Auto Upload Dlls",
		"/v1/module/white/list":    "Download Dll Whitelist",
		"/v1/module/black/list":    "Download Dll Blacklist",

		"/v1/getdirinfo":   "Get Directory Info",
		"/v1/putdir":       "Upload Directory Info",
		"/v1/getdirpolicy": "Download Anti-Tampering Policy",

		"/v1/getprotect": "Download Ransomware Protection Policy",

		"/v1/upserviceport":  "Upload Service Ports",
		"/v1/getvirtualport": "Download Virtual Ports",

		"/v1/getplugging":    "Download Dynamic Blocking Policy",
		"/v1/getipblacklist": "Download IP Blacklist",

		"/v1/getOutreachDetect": "Download Outreach Detection List",

		"/v1/addperipherals":      "Upload Peripheral List",
		"/v1/getwhiteperipherals": "Get Peripheral Whitelist",
		"/v1/getblackperipherals": "Get Peripheral Blacklist",

		"/v1/alertupload": "Upload Alert (except Port)",
		"/v1/upopenport":  "Upload Port Deception Alert",

		"/v1/putsyslog": "Upload Communication Log",

		"/v1/puthardwareinfo": "Upload Agent Resource Info",
		"/v1/uploadIp":        "Upload Agent IP Info",
		"/v1/download":        "Upgrade",
		"/v1/uninstall":       "Uninstall",
		"/v1/error/upload":    "Upload Error File",
		"/v1/gettrustdir":     "Get Trust Directory",

		"/v1/getPwJump":      "Get Password Hopping",
		"/v1/putPwJump":      "Report Password Hopping Result",
		"/v1/getIpJump":      "Get IP Hopping",
		"/v1/getBackups":     "Get Backup Policy",
		"/v1/uploadBackup":   "Upload Backup Result",
		"/v1/getRollbacks":   "Get Rollback Policy",
		"/v1/uploadRollback": "Upload Rollback Result",
	},
}

var WebMapStr = map[string][]string{
	"path": {
		"/v1/auth",
		"/v1/getconf",
		"/v1/gettask",
		"/v1/closetask",
		"/v1/upload/suffix/exe",
		"/v1/uploadproc",
		"/v1/autouploadprocess",
		"/v1/getprocwl",
		"/v1/getprocbl",
		"/v1/module/process",
		"/v1/upload/proc/module",
		"/v1/getdraw",
		"/v1/uploaddraw",
		"/v1/upload/suffix/dll",
		"/v1/pload/allproc/module",
		"/v1/auto/upload/module",
		"/v1/module/white/list",
		"/v1/module/black/list",
		"/v1/getdirinfo",
		"/v1/putdir",
		"/v1/getdirpolicy",
		"/v1/getprotect",
		"/v1/upserviceport",
		"/v1/getvirtualport",
		"/v1/getplugging",
		"/v1/getipblacklist",
		"/v1/addperipherals",
		"/v1/getwhiteperipherals",
		"/v1/getblackperipherals",
		"/v1/getOutreachDetect",
		"/v1/alertupload",
		"/v1/upopenport",
		"/v1/putsyslog",
		"/v1/puthardwareinfo",
		"/v1/uploadIp",
		"/v1/download",
		"/v1/uninstall",
		"/v1/error/upload",
		"/v1/gettrustdir",
		"/v1/getPwJump",
		"/v1/putPwJump",
		"/v1/getIpJump",
		"/v1/getBackups",
		"/v1/uploadBackup",
		"/v1/getRollbacks",
		"/v1/uploadRollback",
	},
}

var EnumMapInt = map[string]map[int]string{
	"status": {
		StatusSuccess: "成功",
		StatusFail:    "异常",
	},
	"status_en": {
		StatusSuccess: "Success",
		StatusFail:    "Abnormal",
	},
}
var WebMapInt = map[string][]int{
	"status": {
		StatusSuccess,
		StatusFail,
	},
}
