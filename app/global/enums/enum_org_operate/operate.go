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
	ModuleOrder                 = "order"
	ModuleAccount               = "account"
	ModuleSyslog                = "syslog"
	ModuleOrg                   = "org"
	ModuleSecurityReport        = "securityReport"
	ModuleOther                 = "other"

	StatusSuccess int = 1
	StatusFail    int = 2

	SourceWeb int = 1
	SourceHs  int = 2
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
		ModuleOrder:          "订单授权",
		ModuleAccount:        "账户管理",
		ModuleOrg:            "企业管理",
		ModuleSyslog:         "通信日志",
		ModuleSecurityReport: "安全报告",
		ModuleOther:          "其他",
	},
}
var WebMapString = map[string][]string{
	"module": {
		ModuleAuth,
		ModuleAgent,
		ModuleProcess,
		ModuleDll,
		ModuleSample,
		ModuleBlackmail,
		ModuleTamper,
		ModulePort,
		ModulePlug,
		ModuleUsb,
		ModuleOutreach,
		ModuleLoginSolve,
		ModuleAlarm,
		ModuleOrder,
		ModuleAccount,
		ModuleOrg,
		ModuleSyslog,
		ModuleSecurityReport,
		ModuleOther,
	},
}

var pathModuleLabel = map[string][]string{
	"/v5/organization/auth/checkUser":     {ModuleAuth, "校验用户"},
	"/v5/organization/auth/getSliderCode": {ModuleAuth, "获取滑块码"},
	"/v5/organization/auth/sendSms":       {ModuleAuth, "发送短信验证码"},
	"/v5/organization/auth/normalLogin":   {ModuleAuth, "密码登录"},
	"/v5/organization/auth/hsLogin":       {ModuleAuth, "幻势登录"},
	//"/v5/organization/auth/logout":        {ModuleAuth, "登出"},

	"/v5/organization/account/saveInfo":    {ModuleAccount, "修改信息"},
	"/v5/organization/account/modifyPw":    {ModuleAccount, "修改密码"},
	"/v5/organization/account/modifyPhone": {ModuleAccount, "修改手机"},
	"/v5/organization/account/modifyEmail": {ModuleAccount, "修改邮箱"},

	"/v5/organization/agent/manage/delete":         {ModuleAgent, "管理页:删除"},
	"/v5/organization/agent/manage/uninstall":      {ModuleAgent, "管理页:卸载"},
	"/v5/organization/agent/manage/upgrade":        {ModuleAgent, "管理页:升级"},
	"/v5/organization/agent/manage/batchDelete":    {ModuleAgent, "管理页:批量删除"},
	"/v5/organization/agent/manage/batchUninstall": {ModuleAgent, "管理页:批量卸载"},
	"/v5/organization/agent/manage/batchUpgrade":   {ModuleAgent, "管理页:批量升级"},
	"/v5/organization/agent/manage/batchModify":    {ModuleAgent, "管理页:批量修改"},
	"/v5/organization/agent/agent/basic/save":      {ModuleAgent, "单页：保存基本信息"},
	"/v5/organization/agent/agent/config/save":     {ModuleAgent, "单页：保存配置信息"},
	"/v5/organization/agent/agent/config/resetUid": {ModuleAgent, "单页：重置终端标识"},
	"/v5/organization/trustDir/global/add":         {ModuleAgent, "全局信任目录：添加"},
	"/v5/organization/trustDir/global/edit":        {ModuleAgent, "全局信任目录：编辑"},
	"/v5/organization/trustDir/global/setSwitch":   {ModuleAgent, "全局信任目录：开关"},
	"/v5/organization/trustDir/global/delete":      {ModuleAgent, "全局信任目录：删除"},

	"/v5/organization/process/manage/setSwitch":        {ModuleProcess, "管理页：开关"},
	"/v5/organization/process/manage/setMode":          {ModuleProcess, "管理页：模式"},
	"/v5/organization/process/manage/batchSetConfig":   {ModuleProcess, "管理页：批量配置"},
	"/v5/organization/process/manage/batchSetWhite":    {ModuleProcess, "管理页：未设置批量白名单"},
	"/v5/organization/process/manage/batchRefresh":     {ModuleProcess, "管理页：批量获刷新"},
	"/v5/organization/process/manage/batchCollect":     {ModuleProcess, "管理页：批量获取所有进程"},
	"/v5/organization/process/agent/setStatus":         {ModuleProcess, "单终端：设置状态"},
	"/v5/organization/process/agent/cancel":            {ModuleProcess, "单终端：取消"},
	"/v5/organization/process/agent/clearAll":          {ModuleProcess, "单终端：清空全部"},
	"/v5/organization/process/agent/clearRunning":      {ModuleProcess, "单终端：清空运行中"},
	"/v5/organization/process/agent/refresh":           {ModuleProcess, "单终端：刷新"},
	"/v5/organization/process/agent/batchWhite":        {ModuleProcess, "单终端：批量白名单"},
	"/v5/organization/process/agent/batchDelete":       {ModuleProcess, "单终端：批量删除"},
	"/v5/organization/process/agent/collect":           {ModuleProcess, "单终端：获取全部进程"},
	"/v5/organization/process/agent/runningDllExtract": {ModuleProcess, "单终端：获取运行中模块"},
	"/v5/organization/process/agent/sampleExtract":     {ModuleProcess, "单终端：提取样本"},
	"/v5/organization/process/agent/setGlobal":         {ModuleProcess, "单终端：设置全局名单"},
	"/v5/organization/process/agent/batchGlobalWhite":  {ModuleProcess, "单终端：批量全局白名单"},
	"/v5/organization/process/global/add":              {ModuleProcess, "全局名单：添加"},
	"/v5/organization/process/global/setSwitch":        {ModuleProcess, "全局名单：开关"},
	"/v5/organization/process/global/delete":           {ModuleProcess, "全局名单：删除"},
	"/v5/organization/process/global/clear":            {ModuleProcess, "全局名单：清空"},

	"/v5/organization/dll/manage/setSwitch":       {ModuleDll, "管理页：开关"},
	"/v5/organization/dll/manage/setMode":         {ModuleDll, "管理页：模式"},
	"/v5/organization/dll/manage/batchSetConfig":  {ModuleDll, "管理页：批量配置"},
	"/v5/organization/dll/manage/batchSetWhite":   {ModuleDll, "管理页：未设置批量白名单"},
	"/v5/organization/dll/manage/batchRefresh":    {ModuleDll, "管理页：批量获刷新"},
	"/v5/organization/dll/manage/batchCollect":    {ModuleDll, "管理页：批量获取所有"},
	"/v5/organization/dll/agent/setStatus":        {ModuleDll, "单终端：设置状态"},
	"/v5/organization/dll/agent/cancel":           {ModuleDll, "单终端：取消"},
	"/v5/organization/dll/agent/clearAll":         {ModuleDll, "单终端：清空全部"},
	"/v5/organization/dll/agent/clearRunning":     {ModuleDll, "单终端：清空运行中"},
	"/v5/organization/dll/agent/refresh":          {ModuleDll, "单终端：刷新"},
	"/v5/organization/dll/agent/batchWhite":       {ModuleDll, "单终端：批量白名单"},
	"/v5/organization/dll/agent/batchDelete":      {ModuleDll, "单终端：批量删除"},
	"/v5/organization/dll/agent/collect":          {ModuleDll, "单终端：获取全部"},
	"/v5/organization/dll/agent/sampleExtract":    {ModuleDll, "单终端：提取样本"},
	"/v5/organization/dll/agent/setGlobal":        {ModuleDll, "单终端：设置全局名单"},
	"/v5/organization/dll/agent/batchGlobalWhite": {ModuleDll, "单终端：批量全局白名单"},
	"/v5/organization/dll/global/add":             {ModuleDll, "全局名单：添加"},
	"/v5/organization/dll/global/setSwitch":       {ModuleDll, "全局名单：开关"},
	"/v5/organization/dll/global/delete":          {ModuleDll, "全局名单：删除"},
	"/v5/organization/dll/global/clear":           {ModuleDll, "全局名单：清空"},

	"/v5/organization/procSample/delete":   {ModuleSample, "删除"},
	"/v5/organization/procSample/clear":    {ModuleSample, "清空"},
	"/v5/organization/procSample/download": {ModuleSample, "下载"},

	"/v5/organization/tamper/manage/setSwitch":      {ModuleTamper, "管理页：开关"},
	"/v5/organization/tamper/manage/setMode":        {ModuleTamper, "管理页：模式"},
	"/v5/organization/tamper/manage/batchSetConfig": {ModuleTamper, "管理页：批量配置"},
	"/v5/organization/tamper/policy/addFolder":      {ModuleTamper, "终端策略：添加目录策略"},
	"/v5/organization/tamper/policy/addFile":        {ModuleTamper, "终端策略：添加文件策略"},
	"/v5/organization/tamper/policy/editFolder":     {ModuleTamper, "终端策略：编辑目录策略"},
	"/v5/organization/tamper/policy/editFile":       {ModuleTamper, "终端策略：编辑文件策略"},
	"/v5/organization/tamper/policy/setSwitch":      {ModuleTamper, "终端策略：开关"},
	"/v5/organization/tamper/policy/delete":         {ModuleTamper, "终端策略：删除"},
	"/v5/organization/tamper/policy/import":         {ModuleTamper, "终端策略：导入"},
	"/v5/organization/tamper/policy/export":         {ModuleTamper, "终端策略：导出"},
	"/v5/organization/tamper/global/addFolder":      {ModuleTamper, "全局策略：添加目录策略"},
	"/v5/organization/tamper/global/addFile":        {ModuleTamper, "全局策略：添加文件策略"},
	"/v5/organization/tamper/global/editFolder":     {ModuleTamper, "全局策略：编辑目录策略"},
	"/v5/organization/tamper/global/editFile":       {ModuleTamper, "全局策略：编辑文件策略"},
	"/v5/organization/tamper/global/setSwitch":      {ModuleTamper, "全局策略：开关"},
	"/v5/organization/tamper/global/delete":         {ModuleTamper, "全局策略：删除"},

	"/v5/organization/blackmail/manage/setSwitch":      {ModuleBlackmail, "管理页：开关"},
	"/v5/organization/blackmail/manage/setMode":        {ModuleBlackmail, "管理页：模式"},
	"/v5/organization/blackmail/manage/batchSetConfig": {ModuleBlackmail, "管理页：批量配置"},
	"/v5/organization/blackmail/policy/add":            {ModuleBlackmail, "终端策略：添加"},
	"/v5/organization/blackmail/policy/edit":           {ModuleBlackmail, "终端策略：编辑"},
	"/v5/organization/blackmail/policy/setSwitch":      {ModuleBlackmail, "终端策略：开关"},
	"/v5/organization/blackmail/policy/delete":         {ModuleBlackmail, "终端策略：删除"},
	"/v5/organization/blackmail/global/add":            {ModuleBlackmail, "全局策略：添加"},
	"/v5/organization/blackmail/global/edit":           {ModuleBlackmail, "全局策略：编辑"},
	"/v5/organization/blackmail/global/setSwitch":      {ModuleBlackmail, "全局策略：开关"},
	"/v5/organization/blackmail/global/delete":         {ModuleBlackmail, "全局策略：删除"},
	"/v5/organization/blackmail/group/add":             {ModuleBlackmail, "文件类型组：添加"},
	"/v5/organization/blackmail/group/edit":            {ModuleBlackmail, "文件类型组：编辑"},
	"/v5/organization/blackmail/group/delete":          {ModuleBlackmail, "文件类型组：删除"},

	"/v5/organization/port/manage/setSwitch":      {ModulePort, "管理页：开关"},
	"/v5/organization/port/manage/batchSetConfig": {ModulePort, "管理页：批量配置"},
	"/v5/organization/port/policy/add":            {ModulePort, "终端策略：添加"},
	"/v5/organization/port/policy/edit":           {ModulePort, "终端策略：编辑"},
	"/v5/organization/port/policy/setSwitch":      {ModulePort, "终端策略：开关"},
	"/v5/organization/port/policy/delete":         {ModulePort, "终端策略：删除"},
	"/v5/organization/port/global/add":            {ModulePort, "全局策略：添加"},
	"/v5/organization/port/global/edit":           {ModulePort, "全局策略：编辑"},
	"/v5/organization/port/global/setSwitch":      {ModulePort, "全局策略：开关"},
	"/v5/organization/port/global/delete":         {ModulePort, "全局策略：删除"},
	"/v5/organization/port/working/refresh":       {ModulePort, "终端正常业务端口：刷新"},

	"/v5/organization/plug/manage/setSwitch":            {ModulePlug, "管理页：开关"},
	"/v5/organization/plug/manage/batchSetConfig":       {ModulePlug, "管理页：批量配置"},
	"/v5/organization/plug/policy/configPort/set":       {ModulePlug, "终端-端口虚开封堵策略：保存"},
	"/v5/organization/plug/policy/configLogin/set":      {ModulePlug, "终端-口令破解封堵策略：保存"},
	"/v5/organization/plug/policy/ipWhite/add":          {ModulePlug, "终端-IP白名单：添加"},
	"/v5/organization/plug/policy/ipWhite/edit":         {ModulePlug, "终端-IP白名单：编辑"},
	"/v5/organization/plug/policy/ipWhite/setSwitch":    {ModulePlug, "终端-IP白名单：开关"},
	"/v5/organization/plug/policy/ipWhite/delete":       {ModulePlug, "终端-IP白名单：删除"},
	"/v5/organization/plug/policy/ipWhite/clear":        {ModulePlug, "终端-IP白名单：清空"},
	"/v5/organization/plug/policy/ipBlack/add":          {ModulePlug, "终端-IP黑名单：添加"},
	"/v5/organization/plug/policy/ipBlack/edit":         {ModulePlug, "终端-IP黑名单：编辑"},
	"/v5/organization/plug/policy/ipBlack/setSwitch":    {ModulePlug, "终端-IP黑名单：开关"},
	"/v5/organization/plug/policy/ipBlack/delete":       {ModulePlug, "终端-IP黑名单：删除"},
	"/v5/organization/plug/policy/ipBlack/clear":        {ModulePlug, "终端-IP黑名单：清空"},
	"/v5/organization/plug/policy/dynamic/disable":      {ModulePlug, "终端-动态封堵：解封"},
	"/v5/organization/plug/policy/dynamic/disableBatch": {ModulePlug, "终端-动态封堵：批量解封"},
	"/v5/organization/plug/policy/dynamic/delete":       {ModulePlug, "终端-动态封堵：批量删除"},
	"/v5/organization/plug/policy/dynamic/deleteBatch":  {ModulePlug, "终端-动态封堵：删除全部"},
	"/v5/organization/plug/global/configPort/set":       {ModulePlug, "全局-端口虚开封堵策略：保存"},
	"/v5/organization/plug/global/configLogin/set":      {ModulePlug, "全局-口令破解封堵策略：保存"},
	"/v5/organization/plug/global/ipWhite/add":          {ModulePlug, "全局-IP白名单：添加"},
	"/v5/organization/plug/global/ipWhite/edit":         {ModulePlug, "全局-IP白名单：编辑"},
	"/v5/organization/plug/global/ipWhite/setSwitch":    {ModulePlug, "全局-IP白名单：开关"},
	"/v5/organization/plug/global/ipWhite/delete":       {ModulePlug, "全局-IP白名单：删除"},
	"/v5/organization/plug/global/ipWhite/clear":        {ModulePlug, "全局-IP白名单：清空"},
	"/v5/organization/plug/global/ipBlack/add":          {ModulePlug, "全局-IP黑名单：添加"},
	"/v5/organization/plug/global/ipBlack/edit":         {ModulePlug, "全局-IP黑名单：编辑"},
	"/v5/organization/plug/global/ipBlack/setSwitch":    {ModulePlug, "全局-IP黑名单：开关"},
	"/v5/organization/plug/global/ipBlack/delete":       {ModulePlug, "全局-IP黑名单：删除"},
	"/v5/organization/plug/global/ipBlack/clear":        {ModulePlug, "全局-IP黑名单：清空"},
	"/v5/organization/plug/global/dynamic/disable":      {ModulePlug, "全局-动态封堵：解封"},
	"/v5/organization/plug/global/dynamic/disableBatch": {ModulePlug, "全局-动态封堵：批量解封"},
	"/v5/organization/plug/global/dynamic/delete":       {ModulePlug, "全局-动态封堵：删除"},
	"/v5/organization/plug/global/dynamic/deleteBatch":  {ModulePlug, "全局-动态封堵：批量删除"},

	"/v5/organization/usb/manage/setSwitch":      {ModuleUsb, "管理页：开关"},
	"/v5/organization/usb/manage/setMode":        {ModuleUsb, "管理页：模式"},
	"/v5/organization/usb/manage/batchSetConfig": {ModuleUsb, "管理页：批量配置"},
	"/v5/organization/usb/policy/setAlias":       {ModuleUsb, "终端策略：设置别名"},
	"/v5/organization/usb/policy/setWhite":       {ModuleUsb, "终端策略：信任"},
	"/v5/organization/usb/policy/setBlack":       {ModuleUsb, "终端策略：禁止"},
	"/v5/organization/usb/policy/refresh":        {ModuleUsb, "终端策略：刷新设备"},
	"/v5/organization/usb/policy/addGlobal":      {ModuleUsb, "终端策略：加入全局名单"},
	"/v5/organization/usb/global/setAlias":       {ModuleUsb, "全局策略：设置别名"},
	"/v5/organization/usb/global/setWhite":       {ModuleUsb, "全局策略：信任"},
	"/v5/organization/usb/global/setBlack":       {ModuleUsb, "全局策略：禁止"},
	"/v5/organization/usb/global/delete":         {ModuleUsb, "全局策略：删除"},

	"/v5/organization/outreach/manage/setSwitch":      {ModuleOutreach, "管理页：开关"},
	"/v5/organization/outreach/manage/batchSetConfig": {ModuleOutreach, "管理页：批量配置"},
	"/v5/organization/outreach/global/add":            {ModuleOutreach, "全局策略：添加"},
	"/v5/organization/outreach/global/fastAdd":        {ModuleOutreach, "全局策略：快捷添加"},
	"/v5/organization/outreach/global/edit":           {ModuleOutreach, "全局策略：编辑"},
	"/v5/organization/outreach/global/setSwitch":      {ModuleOutreach, "全局策略：开关"},
	"/v5/organization/outreach/global/delete":         {ModuleOutreach, "全局策略：删除"},
	"/v5/organization/outreach/detect/add":            {ModuleOutreach, "主动探测：添加"},
	"/v5/organization/outreach/detect/edit":           {ModuleOutreach, "主动探测：编辑"},
	"/v5/organization/outreach/detect/setSwitch":      {ModuleOutreach, "主动探测：开关"},
	"/v5/organization/outreach/detect/delete":         {ModuleOutreach, "主动探测：删除"},

	"/v5/organization/loginSolve/manage/setSwitch":      {ModuleLoginSolve, "管理页：开关"},
	"/v5/organization/loginSolve/manage/batchSetConfig": {ModuleLoginSolve, "管理页：批量配置"},
	"/v5/organization/loginSolve/policyConfig/set":      {ModuleLoginSolve, "终端策略：保存"},
	"/v5/organization/loginSolve/globalConfig/set":      {ModuleLoginSolve, "全局策略：保存"},

	"/v5/organization/alarm/manage/clear":        {ModuleAlarm, "管理页：清空"},
	"/v5/organization/alarm/manage/batchClear":   {ModuleAlarm, "管理页：批量清空"},
	"/v5/organization/alarm/process/handle":      {ModuleAlarm, "进程告警：处理"},
	"/v5/organization/alarm/process/handleAll":   {ModuleAlarm, "进程告警：处理全部"},
	"/v5/organization/alarm/process/clear":       {ModuleAlarm, "进程告警：清空"},
	"/v5/organization/alarm/dll/handle":          {ModuleAlarm, "模块告警：处理"},
	"/v5/organization/alarm/dll/handleAll":       {ModuleAlarm, "模块告警：处理全部"},
	"/v5/organization/alarm/dll/clear":           {ModuleAlarm, "模块告警：清空"},
	"/v5/organization/alarm/tamper/handle":       {ModuleAlarm, "防篡改告警：处理"},
	"/v5/organization/alarm/tamper/handleAll":    {ModuleAlarm, "防篡改告警：处理全部"},
	"/v5/organization/alarm/tamper/clear":        {ModuleAlarm, "防篡改告警：清空"},
	"/v5/organization/alarm/blackmail/handle":    {ModuleAlarm, "勒索防护告警：处理"},
	"/v5/organization/alarm/blackmail/handleAll": {ModuleAlarm, "勒索防护告警：处理全部"},
	"/v5/organization/alarm/blackmail/clear":     {ModuleAlarm, "勒索防护告警：清空"},
	"/v5/organization/alarm/port/handle":         {ModuleAlarm, "端口虚开告警：处理"},
	"/v5/organization/alarm/port/handleAll":      {ModuleAlarm, "端口虚开告警：处理全部"},
	"/v5/organization/alarm/port/clear":          {ModuleAlarm, "端口虚开告警：清空"},
	"/v5/organization/alarm/plug/handle":         {ModuleAlarm, "动态封堵告警：处理"},
	"/v5/organization/alarm/plug/handleAll":      {ModuleAlarm, "动态封堵告警：处理全部"},
	"/v5/organization/alarm/plug/clear":          {ModuleAlarm, "动态封堵告警：清空"},
	"/v5/organization/alarm/usb/handle":          {ModuleAlarm, "外接设备告警：处理"},
	"/v5/organization/alarm/usb/handleAll":       {ModuleAlarm, "外联管控告警：处理全部"},
	"/v5/organization/alarm/usb/clear":           {ModuleAlarm, "外联管控告警：清空"},
	"/v5/organization/alarm/outreach/handle":     {ModuleAlarm, "外联管控告警：处理"},
	"/v5/organization/alarm/outreach/handleAll":  {ModuleAlarm, "外联管控告警：处理全部"},
	"/v5/organization/alarm/outreach/clear":      {ModuleAlarm, "外联管控告警：清空"},
	"/v5/organization/alarm/login/handle":        {ModuleAlarm, "口令破解告警：处理"},
	"/v5/organization/alarm/login/handleAll":     {ModuleAlarm, "口令破解告警：处理全部"},
	"/v5/organization/alarm/login/clear":         {ModuleAlarm, "口令破解告警：清空"},
	"/v5/organization/alarm/self/handle":         {ModuleAlarm, "自保告警：处理"},
	"/v5/organization/alarm/self/handleAll":      {ModuleAlarm, "自保告警：处理全部"},
	"/v5/organization/alarm/self/clear":          {ModuleAlarm, "自保告警：清空"},

	"/v5/organization/syslog/dns/batchDelete":     {ModuleSyslog, "DNS日志：批量删除"},
	"/v5/organization/syslog/dns/clear":           {ModuleSyslog, "DNS日志：清空"},
	"/v5/organization/syslog/inner/batchDelete":   {ModuleSyslog, "对内通讯日志：批量删除"},
	"/v5/organization/syslog/inner/clear":         {ModuleSyslog, "对内通讯日志：清空"},
	"/v5/organization/syslog/outer/batchDelete":   {ModuleSyslog, "对外通讯日志：批量删除"},
	"/v5/organization/syslog/outer/clear":         {ModuleSyslog, "对外通讯日志：清空"},
	"/v5/organization/syslog/process/batchDelete": {ModuleSyslog, "进程启动日志：批量删除"},
	"/v5/organization/syslog/process/clear":       {ModuleSyslog, "进程启动日志：清空"},
	"/v5/organization/syslog/login/batchDelete":   {ModuleSyslog, "登录日志：批量删除"},
	"/v5/organization/syslog/login/clear":         {ModuleSyslog, "登录日志：清空"},

	"/v5/organization/order/purchase/add":           {ModuleOrder, "购买记录：购买"},
	"/v5/organization/order/purchase/cancel":        {ModuleOrder, "购买记录：取消"},
	"/v5/organization/order/purchase/delete":        {ModuleOrder, "购买记录：删除"},
	"/v5/organization/order/purchase/payByBankcard": {ModuleOrder, "购买记录：对公转账方式付款"},
	"/v5/organization/order/grant/save":             {ModuleOrder, "终端授权：保存"},
	"/v5/organization/order/license/import":         {ModuleOrder, "授权许可:导入"},

	"/v5/organization/safetyReport/alarmNotice/save":    {ModuleSecurityReport, "告警通知：保存"},
	"/v5/organization/safetyReport/hardwareNotice/save": {ModuleSecurityReport, "状态：保存"},

	"/v5/organization/department/add":    {ModuleOrg, "组织架构：添加"},
	"/v5/organization/department/edit":   {ModuleOrg, "组织架构：编辑"},
	"/v5/organization/department/delete": {ModuleOrg, "组织架构：删除"},
	"/v5/organization/role/add":          {ModuleOrg, "角色管理：添加"},
	"/v5/organization/role/edit":         {ModuleOrg, "角色管理：编辑"},
	"/v5/organization/role/delete":       {ModuleOrg, "角色管理：删除"},
	"/v5/organization/role/saveMenus":    {ModuleOrg, "角色管理：设置权限"},
	"/v5/organization/user/add":          {ModuleOrg, "子用户：添加"},
	"/v5/organization/user/edit":         {ModuleOrg, "子用户：编辑"},
	"/v5/organization/user/delete":       {ModuleOrg, "子用户：删除"},
	"/v5/organization/user/modifyPw":     {ModuleOrg, "子用户：修改密码"},

	"/v5/organization/task/history/clear": {ModuleOther, "终端任务：清空"},
}

func GetModuleAndLabelByPath(path string) (module string, label string) {
	val, ok := pathModuleLabel[path]
	if ok == false {
		return module, label
	}
	return val[0], val[1]
}
