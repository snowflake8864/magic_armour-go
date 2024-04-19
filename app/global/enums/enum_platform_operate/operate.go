package enum_platform_operate

// 任务

const (
	ModuleAuth = "auth"
	ModuleOrg  = "org"
)

var EnumMap = map[string]map[string]string{
	"module": {
		ModuleAuth: "安全认证",
		ModuleOrg:  "客户管理",
	},
}

var pathModuleLabel = map[string][]string{
	"/platform/auth/checkUser":                    {ModuleAuth, "校验管理员"},
	"/platform/auth/getSliderCode":                {ModuleAuth, "获取滑块码"},
	"/platform/auth/sendSms":                      {ModuleAuth, "发送短信验证码"},
	"/platform/auth/normalLogin":                  {ModuleAuth, "管理员登录"},
	"/platform/auth/platform":                     {ModuleAuth, "管理员登出"},
	"/platform/organize/manage/lock":              {ModuleOrg, "账户管理-锁定"},
	"/platform/organize/manage/unlock":            {ModuleOrg, "账户管理-解锁"},
	"/platform/organize/manage/clean":             {ModuleOrg, "账户管理-清理业务数据"},
	"/platform/organize/manage/setHsSwitch":       {ModuleOrg, "设置幻势开关"},
	"/platform/organize/manage/setLoginSmsSwitch": {ModuleOrg, "设置企业登录短信验证开关"},
	"/platform/organize/verify/add":               {ModuleOrg, "添加企业账户"},
	"/platform/organize/verify/approve":           {ModuleOrg, "账户审核：批准"},
	"/platform/organize/verify/delete":            {ModuleOrg, "账户审核：删除"},
}

func GetModuleAndLabelByPath(path string) (module string, label string) {
	val, ok := pathModuleLabel[path]
	if ok == false {
		return module, label
	}
	return val[0], val[1]
}
