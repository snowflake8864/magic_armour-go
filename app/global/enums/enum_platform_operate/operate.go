package enum_platform_operate

// 任务

const (
	ModuleAuth       = "auth"
	ModuleOrg        = "org"
	ModuleProduct    = "product"
	ModuleOrder      = "order"
	ModuleAgent      = "agent"
	ModuleSoftware   = "software"
	ModuleThreaten   = "threaten"
	ModuleAccount    = "account"
	ModuleDictionary = "dictionary"
	ModuleSystem     = "system"
)

var EnumMap = map[string]map[string]string{
	"module": {
		ModuleAuth:       "安全认证",
		ModuleOrg:        "客户管理",
		ModuleProduct:    "产品管理",
		ModuleOrder:      "订单管理",
		ModuleAgent:      " 终端管理",
		ModuleSoftware:   " 软件管理",
		ModuleThreaten:   "威胁情报",
		ModuleAccount:    "账户管理",
		ModuleDictionary: "数据字典",
		ModuleSystem:     "系统管理",
	},
}

var WebMapString = map[string][]string{
	"module": {
		ModuleAuth,
		ModuleOrg,
		ModuleProduct,
		ModuleOrder,
		ModuleAgent,
		ModuleSoftware,
		ModuleThreaten,
		ModuleAccount,
		ModuleDictionary,
		ModuleSystem,
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
	"/platform/product/manage/add":                {ModuleProduct, "添加"},
	"/platform/product/manage/edit":               {ModuleProduct, "编辑"},
	"/platform/product/manage/delete":             {ModuleProduct, "删除"},
	"/platform/agent/org/addGrant":                {ModuleAgent, "添加授权"},
	"/platform/order/manage/decline":              {ModuleOrder, "直降"},
	"/platform/order/manage/cancel":               {ModuleOrder, "取消"},
	"/platform/order/manage/paid":                 {ModuleOrder, "确认已付款"},
	"/platform/order/manage/unpaid":               {ModuleOrder, "确认未付款"},
	"/platform/software/manage/add":               {ModuleSoftware, "添加"},
	"/platform/software/manage/delete":            {ModuleSoftware, "删除"},
	"/platform/software/manage/publish":           {ModuleSoftware, "发布"},
	"/platform/software/manage/cancel":            {ModuleSoftware, "取消"},
	"/platform/software/manage/clear":             {ModuleSoftware, "清空"},
	"/platform/threaten/cache/delete":             {ModuleThreaten, "缓存-删除"},
	"/platform/threaten/cache/clear":              {ModuleThreaten, "缓存-清空"},
	"/platform/threaten/cache/refresh":            {ModuleThreaten, "缓存-刷新"},
	"/platform/threaten/cache/batchWhite":         {ModuleThreaten, "缓存-批量白名单"},
	"/platform/threaten/cache/batchBlack":         {ModuleThreaten, "缓存-批量黑名单"},
	"/platform/threaten/push/clear":               {ModuleThreaten, "推送-清空"},
	"/platform/account/modifyPw":                  {ModuleAccount, "修改密码"},
	"/platform/dictionary/save":                   {ModuleDictionary, "修改"},
	"/platform/system/agentLogConfig/set":         {ModuleSystem, "修改终端日志配置"},
	"/platform/system/network/set":                {ModuleSystem, "修改网络配置"},
	"/platform/system/ssh/set":                    {ModuleSystem, "修改SSH开关"},
	"/platform/system/time/set":                   {ModuleSystem, "修改系统时间"},
}

func GetModuleAndLabelByPath(path string) (module string, label string) {
	val, ok := pathModuleLabel[path]
	if ok == false {
		return module, label
	}
	return val[0], val[1]
}
