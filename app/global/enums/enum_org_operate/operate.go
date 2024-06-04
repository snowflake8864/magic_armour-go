package enum_org_operate

// 任务

const (
	ModuleAuth           string = "auth"
	ModuleRegistry              = "registry"
	ModuleFindPw                = "findPw"
	ModuleAgent                 = "agent"
	ModuleProcess               = "process"
	ModuleDll                   = "dll"
	ModuleSample                = "sample"
	ModuleBlackmail             = "blackmail"
	ModuleTamper                = "tamper"
	ModulePort                  = "port"
	ModulePlug                  = "plug"
	ModuleUsb                   = "usb"
	ModuleOutreach              = "outreach"
	ModuleLoginSolve            = "loginSolve"
	ModuleAlarm                 = "alarm"
	ModuleGrant                 = "grant"
	ModuleAccount               = "account"
	ModuleSyslog                = "syslog"
	ModuleOrg                   = "org"
	ModuleSecurityReport        = "securityReport"
	ModuleOther                 = "other"
)

var EnumMap = map[string]map[string]string{
	"module": {
		ModuleAuth:           "安全认证",
		ModuleAgent:          "终端管理",
		ModuleProcess:        "可信防护-进程管理",
		ModuleDll:            "可信防护-模块管理",
		ModuleSample:         "可信防护-样本管理",
		ModuleBlackmail:      "勒索防护",
		ModuleTamper:         "防篡改",
		ModulePort:           "端口虚开",
		ModulePlug:           "动态封堵",
		ModuleUsb:            "外接设备",
		ModuleOutreach:       "外联管控",
		ModuleLoginSolve:     "口令破解",
		ModuleAlarm:          "告警管理",
		ModuleGrant:          "授权管理",
		ModuleAccount:        "账户管理",
		ModuleOrg:            "企业管理",
		ModuleSyslog:         "通信日志",
		ModuleSecurityReport: "安全报告",
		ModuleOther:          "其他",
	},
}

var pathModuleLabel = map[string][]string{
	"/v5/organization/auth/checkUser":                {ModuleAuth, "校验用户"},
	"/v5/organization/auth/getSliderCode":            {ModuleAuth, "获取滑块码"},
	"/v5/organization/auth/sendSms":                  {ModuleAuth, "发送短信验证码"},
	"/v5/organization/auth/normalLogin":              {ModuleAuth, "密码登录"},
	"/v5/organization/auth/hsLogin":                  {ModuleAuth, "幻势登录"},
	"/v5/organization/auth/logout":                   {ModuleAuth, "幻势登录"},
	"/v5/organization/agent/manage/delete":           {ModuleAgent, "管理页:删除"},
	"/v5/organization/agent/manage/uninstall":        {ModuleAgent, "管理页:卸载"},
	"/v5/organization/agent/manage/upgrade":          {ModuleAgent, "管理页:卸载"},
	"/v5/organization/agent/manage/batchDelete":      {ModuleAgent, "管理页:批量删除"},
	"/v5/organization/agent/manage/batchUninstall":   {ModuleAgent, "管理页:批量卸载"},
	"/v5/organization/agent/manage/batchUpgrade":     {ModuleAgent, "管理页:批量升级"},
	"/v5/organization/agent/manage/batchModify":      {ModuleAgent, "管理页:批量修改"},
	"/v5/organization/agent/agent/basic/save":        {ModuleAgent, "单页：保存基本信息"},
	"/v5/organization/agent/agent/config/save":       {ModuleAgent, "单页：保存配置信息"},
	"/v5/organization/process/manage/setSwitch":      {ModuleProcess, "管理页：开关"},
	"/v5/organization/process/manage/setMode":        {ModuleProcess, "管理页：模式"},
	"/v5/organization/process/manage/batchSetConfig": {ModuleProcess, "管理页：批量设置配置"},
	"/v5/organization/process/manage/batchSetWhite":  {ModuleProcess, "管理页：未设置批量白名单"},
	"/v5/organization/process/manage/batchRefresh":   {ModuleProcess, "管理页：批量获刷新"},
	"/v5/organization/process/manage/batchCollect":   {ModuleProcess, "管理页：批量获取所有进程"},
	"/v5/organization/process/agent/setStatus":       {ModuleProcess, "终端页：设置状态"},
}

func GetModuleAndLabelByPath(path string) (module string, label string) {
	val, ok := pathModuleLabel[path]
	if ok == false {
		return module, label
	}
	return val[0], val[1]
}
