package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/response"
)

// AgentController 终端相关
type AgentController struct {
}

// GetConfig 获取配置信息
func (u *AgentController) GetConfig(ctx *gin.Context) {
	//关闭任务
	currentAgent := cur_userinfo.GetCurrentAgent(ctx)
	go task_service.CreateCurdFactory().CloseTasks(currentAgent.AgentId, []int{enum_task.TypeAgentConfigDownload}, ctx.ClientIP())

	//终端配置
	agentCurdFact := agent_service.CreateCurdFactory()
	agentModel := agentCurdFact.GetItemByAgent(currentAgent.AgentId, agentCurdFact.ColumnConfig(true))

	//数据字典
	dictIds := []string{
		enum_dictionary.AgentLogSendProto,
		enum_dictionary.SystemAgentProto,
		enum_dictionary.SystemAgentIp,
		enum_dictionary.SystemAgentPort,
	}
	mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)

	mapConf := make(map[string]interface{})
	//服务器地址
	mapConf["serveripport"] = mapDict[enum_dictionary.SystemAgentProto] + "://" + mapDict[enum_dictionary.SystemAgentIp] + ":" + mapDict[enum_dictionary.SystemAgentPort]
	//暂时无使用
	mapConf["logipport"] = ""
	mapConf["logproto"] = helpers.String2Int(mapDict[enum_dictionary.AgentLogSendProto])
	mapConf["logsent"] = 0
	mapConf["api_port"] = helpers.String2Int(mapDict[enum_dictionary.SystemAgentPort])
	mapConf["syslog_port"] = 0

	//进程
	mapConf["proc_switch"] = agentModel.ProcessSwitch
	mapConf["proc_protect"] = agentModel.ProcessMode
	//模块
	mapConf["module_switch"] = agentModel.DllSwitch
	mapConf["module_protect"] = agentModel.DllMode
	//防篡改
	mapConf["file_switch"] = agentModel.TamperSwitch
	mapConf["file_protect"] = agentModel.TamperMode
	//勒索
	mapConf["extortion_switch"] = agentModel.BlackmailSwitch
	mapConf["extortion_protect"] = agentModel.BlackmailMode
	//端口虚开
	mapConf["open_port_switch"] = agentModel.PortSwitch
	//封堵
	mapConf["dynamic_switch"] = agentModel.PlugSwitch
	//外设
	mapConf["usb_switch"] = agentModel.UsbSwitch
	mapConf["usb_protect"] = agentModel.UsbMode
	//自保
	mapConf["self_protect_switch"] = agentModel.SelfProtectSwitch
	//内核升级
	mapConf["kernel_upgrade_switch"] = agentModel.KernelUpgradeSwitch
	//agent通讯日志
	mapConf["syslog_switch"] = enums.SwitchOff
	mapConf["syslog_dns_switch"] = agentModel.SyslogDnsSwitch
	mapConf["syslog_outer_switch"] = agentModel.SyslogOuterSwitch
	mapConf["syslog_inner_switch"] = agentModel.SyslogInnerSwitch
	mapConf["syslog_process_switch"] = agentModel.SyslogProcessSwitch
	mapConf["syslog_login_switch"] = agentModel.SyslogLoginSwitch
	//上传硬件信息间隔时间(秒)
	mapConf["hardware_switch"] = agentModel.HardwareSwitch
	mapConf["hardware_time"] = agentModel.HardwareTime

	//linux不支持的功能
	if currentAgent.PlatformType == enum_software.PlatformTypeLinux {
		mapConf["module_switch"] = enums.SwitchOff
		mapConf["module_protect"] = enum_agent.ModeMonitor
	}
	//win xp 不支持的功能
	if currentAgent.OsType == enum_software.OsTypeWinXp {
		mapConf["open_port_switch"] = enums.SwitchOff
		mapConf["dynamic_switch"] = enum_agent.ModeMonitor
	}

	//心跳间隔时间(秒)
	if heartFlag := task_service.CreateCacheFactory().GetHeartbeatFlag(currentAgent.AgentId); heartFlag {
		mapConf["crontime"] = agentModel.CommonTime
	} else {
		mapConf["crontime"] = agentModel.FastTime
	}

	response.SuccessReturnData(ctx, gin.H{"conf": mapConf})
}

// GetInfo 获取终端信息(调试用)
func (u *AgentController) GetInfo(ctx *gin.Context) {
	currentAgent := cur_userinfo.GetCurrentAgent(ctx)
	response.SuccessReturnData(ctx, currentAgent)
}

// Upgrade 升级
func (u *AgentController) Upgrade(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "升级")
}

// Uninstall 卸载回调
func (u *AgentController) Uninstall(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "卸载回调")
}

// GetTrustDir 获取信任目录
func (u *AgentController) GetTrustDir(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "获取信任目录")
}

// UploadErrorFile 上传转储(错误) 文件
func (u *AgentController) UploadErrorFile(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传转储(错误) 文件")
}

// MonitorHardware 监控硬件资源
func (u *AgentController) MonitorHardware(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "监控硬件资源")
}
