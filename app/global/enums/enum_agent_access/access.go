package enum_agent_access

// 终端

const ()

var EnumStrMap = map[string]map[string]string{
	"path": {
		"v1/auth":      "登录",
		"v1/getconf":   "下载配置信息",
		"v1/gettask":   "获取任务(心跳)",
		"v1/closetask": "关闭任务",

		"v1/upload/suffix/exe":    "上传全部程序",
		"v1/uploadproc":           "上传运行中进程",
		"v1/autouploadprocess":    "主动上传进程",
		"v1/getprocwl":            "下载进程白名单",
		"v1/getprocbl":            "下载进程黑名单",
		"v1/module/process":       "进程对应模块:下载指定进程",
		"v1/upload/proc/module":   "进程对应模块:上传指定进程下的模块",
		"v1/getdraw":              "查询未提取样本",
		"v1/uploaddraw":           "上传样本",
		"v1/upload/suffix/dll":    "上传全部模块",
		"v1/pload/allproc/module": "上传运行中模块",
		"v1/auto/upload/module":   "主动上传模块",
		"v1/module/white/list":    "下载模块白名单",
		"v1/module/black/list":    "下载模块黑名单",

		"v1/getdirinfo":   "获取查询目录",
		"v1/putdir":       "上传目录信息",
		"v1/getdirpolicy": "下载防篡改策略",

		"v1/getprotect": "下载勒索防护策略",

		"v1/upserviceport":  "上传业务端口",
		"v1/getvirtualport": "下载虚开端口",

		"v1/getplugging":    "下载动态封堵策略",
		"v1/getipblacklist": "下载封堵ip黑名单",

		"v1/addperipherals":      "上传外设列表",
		"v1/getwhiteperipherals": "获取外设信任（白名单）列表",
		"v1/getblackperipherals": "获取外设禁止（黑名单）列表",

		"v1/alertupload": "上传告警(除端口端口外)",
		"v1/upopenport":  "上传端口虚开告警",

		"v1/putsyslog": "上传通信日志",

		"v1/puthardwareinfo": "上传终端资源信息",
		"v1/download":        "升级",
		"v1/uninstall":       "卸载",
		"v1/error/upload":    "上传错误文件",
		"v1/gettrustdir":     "获取信任目录",
	},
}
