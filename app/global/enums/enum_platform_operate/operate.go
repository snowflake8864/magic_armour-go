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

	StatusSuccess int = 1
	StatusFail    int = 2
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
	"module_en": {
		ModuleAuth:       "Authentication",
		ModuleOrg:        "Customer Management",
		ModuleProduct:    "Product Management",
		ModuleOrder:      "Order Management",
		ModuleAgent:      "Agent Management",
		ModuleSoftware:   "Software Management",
		ModuleThreaten:   "Threat Intelligence",
		ModuleAccount:    "Account Management",
		ModuleDictionary: "Data Dictionary",
		ModuleSystem:     "System Management",
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
	"/v5/platform/auth/checkUser":     {ModuleAuth, "校验管理员", "Validate Admin"},
	"/v5/platform/auth/getSliderCode": {ModuleAuth, "获取滑块码", "Get Slider Code"},
	"/v5/platform/auth/sendSms":       {ModuleAuth, "发送短信验证码", "Send SMS Code"},
	"/v5/platform/auth/normalLogin":   {ModuleAuth, "管理员登录", "Admin Login"},
	//"/v5/platform/auth/logout":        {ModuleAuth, "管理员登出", "Admin Logout"},

	"/v5/platform/organize/manage/lock":              {ModuleOrg, "账户管理-锁定", "Account Management - Lock"},
	"/v5/platform/organize/manage/unlock":            {ModuleOrg, "账户管理-解锁", "Account Management - Unlock"},
	"/v5/platform/organize/manage/clean":             {ModuleOrg, "账户管理-清理业务数据", "Account Management - Clear Business Data"},
	"/v5/platform/organize/manage/setHsSwitch":       {ModuleOrg, "设置幻势开关", "Set Huanshi Switch"},
	"/v5/platform/organize/manage/setLoginSmsSwitch": {ModuleOrg, "设置企业登录短信验证开关", "Set Enterprise Login SMS Verification Switch"},
	"/v5/platform/organize/verify/add":               {ModuleOrg, "添加企业账户", "Add Enterprise Account"},
	"/v5/platform/organize/verify/approve":           {ModuleOrg, "账户审核：批准", "Account Approval: Approve"},
	"/v5/platform/organize/verify/delete":            {ModuleOrg, "账户审核：删除", "Account Approval: Delete"},
	"/v5/platform/product/manage/add":                {ModuleProduct, "添加", "Add"},
	"/v5/platform/product/manage/edit":               {ModuleProduct, "编辑", "Edit"},
	"/v5/platform/product/manage/delete":             {ModuleProduct, "删除", "Delete"},
	"/v5/platform/agent/org/addGrant":                {ModuleAgent, "添加授权", "Add Grant"},
	"/v5/platform/order/manage/decline":              {ModuleOrder, "直降", "Direct Discount"},
	"/v5/platform/order/manage/cancel":               {ModuleOrder, "取消", "Cancel"},
	"/v5/platform/order/manage/paid":                 {ModuleOrder, "确认已付款", "Confirm Paid"},
	"/v5/platform/order/manage/unpaid":               {ModuleOrder, "确认未付款", "Confirm Unpaid"},
	"/v5/platform/software/manage/add":               {ModuleSoftware, "添加", "Add"},
	"/v5/platform/software/manage/delete":            {ModuleSoftware, "删除", "Delete"},
	"/v5/platform/software/manage/publish":           {ModuleSoftware, "发布", "Publish"},
	"/v5/platform/software/manage/cancel":            {ModuleSoftware, "取消", "Cancel"},
	"/v5/platform/software/manage/clear":             {ModuleSoftware, "清空", "Clear"},
	"/v5/platform/software/agentConfig/set":          {ModuleSoftware, "终端配置", "Agent Config"},
	"/v5/platform/threaten/config/online":            {ModuleThreaten, "配置-验证", "Config - Validate"},
	"/v5/platform/threaten/config/save":              {ModuleThreaten, "配置-保存", "Config - Save"},
	"/v5/platform/threaten/cache/delete":             {ModuleThreaten, "缓存-删除", "Cache - Delete"},
	"/v5/platform/threaten/cache/clear":              {ModuleThreaten, "缓存-清空", "Cache - Clear"},
	"/v5/platform/threaten/cache/refresh":            {ModuleThreaten, "缓存-刷新", "Cache - Refresh"},
	"/v5/platform/threaten/cache/import":             {ModuleThreaten, "缓存-导入", "Cache - Import"},
	"/v5/platform/threaten/cache/batchWhite":         {ModuleThreaten, "缓存-批量白名单", "Cache - Batch Whitelist"},
	"/v5/platform/threaten/cache/batchBlack":         {ModuleThreaten, "缓存-批量黑名单", "Cache - Batch Blacklist"},
	"/v5/platform/threaten/push/clear":               {ModuleThreaten, "推送-清空", "Push - Clear"},
	"/v5/platform/account/modifyPw":                  {ModuleAccount, "修改密码", "Modify Password"},
	"/v5/platform/dictionary/save":                   {ModuleDictionary, "修改", "Modify"},
	"/v5/platform/system/agentLogConfig/set":         {ModuleSystem, "修改终端日志配置", "Modify Agent Log Config"},
	"/v5/platform/system/network/set":                {ModuleSystem, "修改网络配置", "Modify Network Config"},
	"/v5/platform/system/ssh/set":                    {ModuleSystem, "修改SSH开关", "Modify SSH Switch"},
	"/v5/platform/system/sntp/set":                   {ModuleSystem, "设置SNTP服务", "Set SNTP Service"},
	"/v5/platform/system/time/set":                   {ModuleSystem, "修改系统时间", "Modify System Time"},
	"/v5/platform/system/upgrade/run":                {ModuleSystem, "系统升级", "System Upgrade"},
	"/v5/platform/system/upgrade/saveLog":            {ModuleSystem, "系统升级", "System Upgrade"},
	"/v5/platform/system/huanshi/setThird":           {ModuleSystem, "设置第三方日志平台", "Set Third-party Log Platform"},
	"/v5/platform/system/huanshi/set":                {ModuleSystem, "设置日志平台", "Set Log Platform"},
}

func GetModuleAndLabelByPath(path string) (module string, label, labelEn string) {
	val, ok := pathModuleLabel[path]
	if ok == false {
		return module, label, labelEn
	}
	return val[0], val[1], val[2]
}
