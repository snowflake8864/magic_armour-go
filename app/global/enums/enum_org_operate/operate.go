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
	"/v5/organization/auth/checkUser":                  {ModuleAuth, "校验用户"},
	"/v5/organization/auth/getSliderCode":              {ModuleAuth, "获取滑块码"},
	"/v5/organization/auth/sendSms":                    {ModuleAuth, "发送短信验证码"},
	"/v5/organization/auth/normalLogin":                {ModuleAuth, "密码登录"},
	"/v5/organization/auth/hsLogin":                    {ModuleAuth, "幻势登录"},
	"/v5/organization/auth/logout":                     {ModuleAuth, "幻势登录"},
	"/v5/organization/agent/manage/delete":             {ModuleAgent, "管理页:删除"},
	"/v5/organization/agent/manage/uninstall":          {ModuleAgent, "管理页:卸载"},
	"/v5/organization/agent/manage/upgrade":            {ModuleAgent, "管理页:卸载"},
	"/v5/organization/agent/manage/batchDelete":        {ModuleAgent, "管理页:批量删除"},
	"/v5/organization/agent/manage/batchUninstall":     {ModuleAgent, "管理页:批量卸载"},
	"/v5/organization/agent/manage/batchUpgrade":       {ModuleAgent, "管理页:批量升级"},
	"/v5/organization/agent/manage/batchModify":        {ModuleAgent, "管理页:批量修改"},
	"/v5/organization/agent/agent/basic/save":          {ModuleAgent, "单页：保存基本信息"},
	"/v5/organization/agent/agent/config/save":         {ModuleAgent, "单页：保存配置信息"},
	"/v5/organization/trustDir/global/add":             {ModuleAgent, "全局信任目录：添加"},
	"/v5/organization/trustDir/global/edit":            {ModuleAgent, "全局信任目录：编辑"},
	"/v5/organization/trustDir/global/setSwitch":       {ModuleAgent, "全局信任目录：开关"},
	"/v5/organization/trustDir/global/delete":          {ModuleAgent, "全局信任目录：删除"},
	"/v5/organization/process/manage/setSwitch":        {ModuleProcess, "管理页：开关"},
	"/v5/organization/process/manage/setMode":          {ModuleProcess, "管理页：模式"},
	"/v5/organization/process/manage/batchSetConfig":   {ModuleProcess, "管理页：批量设置配置"},
	"/v5/organization/process/manage/batchSetWhite":    {ModuleProcess, "管理页：未设置批量白名单"},
	"/v5/organization/process/manage/batchRefresh":     {ModuleProcess, "管理页：批量获刷新"},
	"/v5/organization/process/manage/batchCollect":     {ModuleProcess, "管理页：批量获取所有进程"},
	"/v5/organization/process/agent/setStatus":         {ModuleProcess, "单终端：设置状态"},
	"/v5/organization/process/agent/cancel":            {ModuleProcess, "单终端：取消"},
	"/v5/organization/process/agent/clearAll":          {ModuleProcess, "单终端：清空全部"},
	"/v5/organization/process/agent/clearRunning":      {ModuleProcess, "单终端：清空运行中"},
	"/v5/organization/process/agent/refresh":           {ModuleProcess, "单终端：刷新"},
	"/v5/organization/process/agent/collect":           {ModuleProcess, "单终端：获取全部进程"},
	"/v5/organization/process/agent/runningDllExtract": {ModuleProcess, "单终端：获取运行中模块"},
	"/v5/organization/process/agent/sampleExtract":     {ModuleProcess, "单终端：提取样本"},
	"/v5/organization/process/agent/setGlobal":         {ModuleProcess, "单终端：设置全局名单"},
	"/v5/organization/process/agent/batchGlobalWhite":  {ModuleProcess, "单终端：批量全局白名单"},
	"/v5/organization/process/global/add":              {ModuleProcess, "全局名单：添加"},
	"/v5/organization/process/global/setSwitch":        {ModuleProcess, "全局名单：开关"},
	"/v5/organization/process/global/delete":           {ModuleProcess, "全局名单：删除"},
	"/v5/organization/process/global/clear":            {ModuleProcess, "全局名单：清空"},
	"/v5/organization/dll/manage/setSwitch":            {ModuleDll, "管理页：开关"},
	"/v5/organization/dll/manage/setMode":              {ModuleDll, "管理页：模式"},
	"/v5/organization/dll/manage/batchSetConfig":       {ModuleDll, "管理页：批量设置配置"},
	"/v5/organization/dll/manage/batchSetWhite":        {ModuleDll, "管理页：未设置批量白名单"},
	"/v5/organization/dll/manage/batchRefresh":         {ModuleDll, "管理页：批量获刷新"},
	"/v5/organization/dll/manage/batchCollect":         {ModuleDll, "管理页：批量获取所有"},
	"/v5/organization/dll/agent/setStatus":             {ModuleDll, "单终端：设置状态"},
	"/v5/organization/dll/agent/cancel":                {ModuleDll, "单终端：取消"},
	"/v5/organization/dll/agent/clearAll":              {ModuleDll, "单终端：清空全部"},
	"/v5/organization/dll/agent/clearRunning":          {ModuleDll, "单终端：清空运行中"},
	"/v5/organization/dll/agent/refresh":               {ModuleDll, "单终端：刷新"},
	"/v5/organization/dll/agent/collect":               {ModuleDll, "单终端：获取全部"},
	"/v5/organization/dll/agent/sampleExtract":         {ModuleDll, "单终端：提取样本"},
	"/v5/organization/dll/agent/setGlobal":             {ModuleDll, "单终端：设置全局名单"},
	"/v5/organization/dll/agent/batchGlobalWhite":      {ModuleDll, "单终端：批量全局白名单"},
	"/v5/organization/dll/global/add":                  {ModuleDll, "全局名单：添加"},
	"/v5/organization/dll/global/setSwitch":            {ModuleDll, "全局名单：开关"},
	"/v5/organization/dll/global/delete":               {ModuleDll, "全局名单：删除"},
	"/v5/organization/dll/global/clear":                {ModuleDll, "全局名单：清空"},
	"/v5/organization/procSample/delete":               {ModuleSample, "删除"},
	"/v5/organization/procSample/clear":                {ModuleSample, "清空"},
	"/v5/organization/procSample/download":             {ModuleSample, "下载"},
	"/v5/organization/tamper/manage/setSwitch":         {ModuleTamper, "管理页：开关"},
	"/v5/organization/tamper/manage/setMode":           {ModuleTamper, "管理页：模式"},
	"/v5/organization/tamper/manage/batchSetConfig":    {ModuleTamper, "管理页：批量设置配置"},
	"/v5/organization/tamper/policy/addFolder":         {ModuleTamper, "策略页：添加文件夹策略"},
}

func GetModuleAndLabelByPath(path string) (module string, label string) {
	val, ok := pathModuleLabel[path]
	if ok == false {
		return module, label
	}
	return val[0], val[1]
}
