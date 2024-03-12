package enum_platform_operate

// 任务

const (
	ModuleAuth             = "auth"
	ModuleAgent            = "agent"
	ModuleProcess          = "process"
	moduleProcessGlobal    = "processGlobal"
	moduleDll              = "dll"
	moduleDllGlobal        = "dllGlobal"
	moduleProcSample       = "procSample"
	moduleBlackmail        = "blackmail"
	moduleBlackmailGroup   = "blackmailGroup"
	moduleTamper           = "tamper"
	moduleTamperGlobal     = "tamperGlobal"
	modulePort             = "port"
	modulePortGlobal       = "portGlobal"
	modulePlugDynamic      = "plugDynamic"
	modulePlugGlobalConfig = "plugGlobalConfig"
	modulePlugGlobalIp     = "plugGlobalIP"
	modulePlugList         = "plugGlobalList"
	moduleUsb              = "usb"
	moduleUsbGlobal        = "usbGlobal"
	moduleOutreach         = "outreach"
	moduleOutreachGlobal   = "outreachGlobal"
	moduleLoginSolve       = "loginSolve"
	moduleLoginSolveGlobal = "loginSolveGlobal"
	moduleAlarm            = "alarm"
	moduleThreaten         = "threaten"
	moduleGrant            = "grant"
	moduleAdmin            = "admin"
	moduleSystem           = "system"
	moduleSyslog           = "syslog"
	moduleOther            = "other"
)

var EnumMap = map[string]map[string]string{
	"module": {
		ModuleAuth:             "安全认证",
		ModuleAgent:            "终端管理",
		ModuleProcess:          "可信防护-进程管理",
		moduleProcessGlobal:    "可信防护-进程全局名单",
		moduleDll:              "可信防护-模块管理",
		moduleDllGlobal:        "可信防护-模块全局名单",
		moduleProcSample:       "可信防护-样本提取",
		moduleBlackmail:        "勒索防护",
		moduleBlackmailGroup:   "勒索防护-文件类型组设置",
		moduleTamper:           "防篡改",
		moduleTamperGlobal:     "全局防篡改",
		modulePort:             "端口虚开",
		modulePortGlobal:       "全局端口虚开",
		modulePlugDynamic:      "动态封堵",
		modulePlugGlobalConfig: "全局封堵配置",
		modulePlugGlobalIp:     "封堵全局名单",
		modulePlugList:         "封堵列表",
		moduleUsb:              "外设管理",
		moduleUsbGlobal:        "外设全局名单",
		moduleOutreach:         "外联管控",
		moduleOutreachGlobal:   "外联白名单",
		moduleLoginSolve:       "口令破解",
		moduleLoginSolveGlobal: "全局口令破解",
		moduleAlarm:            "告警",
		moduleThreaten:         "威胁情报",
		moduleGrant:            "授权管理",
		moduleAdmin:            "平用户管理",
		moduleSystem:           "系统管理",
		moduleSyslog:           "通信日志",
		moduleOther:            "其他",
	},
}

var pathModuleLabel = map[string][]string{
	"/platform/auth/checkUser":     {ModuleAuth, "校验用户"},
	"/platform/auth/getSliderCode": {ModuleAuth, "获取滑块码"},
	"/platform/auth/sendSms":       {ModuleAuth, "发送短信验证码"},
	"/platform/auth/normalLogin":   {ModuleAuth, "密码登录"},
	"/platform/auth/hsLogin":       {ModuleAuth, "幻势登录"},
}

func GetModuleAndLabelByPath(path string) (module string, label string) {
	val, ok := pathModuleLabel[path]
	if ok == false {
		return module, label
	}
	return val[0], val[1]
}
