package agent_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_grant"
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/service/platform/department_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/grant_service"
	"goskeleton/app/utils/helpers"
)

func CreateAuthFactory() *AgentAuth {
	return &AgentAuth{masterModel: model_agent.CreateAgentFactory()}
}

type AgentAuth struct {
	AgentBase
	masterModel *model_agent.AgentModel
}

// RegisterOrLogin 存在登录，不存在注册
func (crd *AgentAuth) RegisterOrLogin(param *model_agent.AgentModel) (bool, string) {
	column := "agent_id,dept_id,platform_type,os_type,process_switch,dll_switch,tamper_switch,blackmail_switch,port_switch,plug_switch,usb_switch"
	agentModel := CreateCurdFactory().GetItemByAgent(param.AgentId, column)
	nowDate := helpers.GetDatetime()
	isAllocateDept := true //是否分配组织架构

	if helpers.StringIsEmpty(agentModel.AgentId) {
		//新终端默认使用全局配置
		dictIds := []string{
			enum_dictionary.AgentGlobalPlugPortConfig,
			enum_dictionary.AgentGlobalPlugLoginConfig,
			enum_dictionary.AgentGlobalLoginSolveConfig,
		}
		mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
		var isExist bool
		jsonPlugPort, isExist := mapDict[enum_dictionary.AgentGlobalPlugPortConfig]
		if !isExist || helpers.StringIsEmpty(jsonPlugPort) {
			plugPortConfig := enum_plug.PlugPortConfig{}
			plugPortConfig.Format()
			jsonPlugPort = helpers.JsonEncode(plugPortConfig)
		}
		jsonPlugLogin, isExist := mapDict[enum_dictionary.AgentGlobalPlugLoginConfig]
		if !isExist || helpers.StringIsEmpty(jsonPlugLogin) {
			plugLoginConfig := enum_plug.PlugLoginConfig{}
			plugLoginConfig.Format()
			jsonPlugLogin = helpers.JsonEncode(plugLoginConfig)
		}
		jsonLoginSolve, isExist := mapDict[enum_dictionary.AgentGlobalPlugLoginConfig]
		if !isExist || helpers.StringIsEmpty(jsonPlugLogin) {
			loginSolveConfig := enum_login_solve.LoginSolveConfig{}
			loginSolveConfig.Format()
			jsonLoginSolve = helpers.JsonEncode(loginSolveConfig)
		}

		//授权配置
		grantConfig := enum_grant.AgentGrantConfig{}
		grantConfig.Default()
		jsonGrant := helpers.JsonEncode(grantConfig)

		param.IsOnline = enums.IsYes
		param.OnlineAt = nowDate
		param.LiveAt = nowDate
		param.IsSafe = enums.IsYes
		param.UpgradeVer = ""
		param.DeptId = 0
		param.AgentType = 0
		param.Alias = ""
		param.ContactName = ""
		param.ContactTel = ""
		param.Position = ""
		param.Remark = ""
		param.GrantConfig = jsonGrant
		param.ProcessSwitch = enums.SwitchOff
		param.ProcessMode = enum_agent.ModeMonitor
		param.DllSwitch = enums.SwitchOff
		param.DllMode = enum_agent.ModeMonitor
		param.TamperSwitch = enums.SwitchOff
		param.TamperMode = enum_agent.ModeMonitor
		param.BlackmailSwitch = enums.SwitchOff
		param.BlackmailMode = enum_agent.ModeMonitor
		param.PortSwitch = enums.SwitchOff
		param.PlugSwitch = enums.SwitchOff
		param.PlugPortConfig = jsonPlugPort
		param.PlugLoginConfig = jsonPlugLogin
		param.UsbSwitch = enums.SwitchOff
		param.UsbMode = enum_agent.ModeMonitor
		param.SyslogDnsSwitch = enums.SwitchOff
		param.SyslogOuterSwitch = enums.SwitchOff
		param.SyslogInnerSwitch = enums.SwitchOff
		param.SyslogProcessSwitch = enums.SwitchOff
		param.SyslogLoginSwitch = enums.SwitchOn
		param.OutreachAlarmSwitch = enums.SwitchOff
		param.SelfProtectSwitch = enums.SwitchOn
		param.KernelUpgradeSwitch = enums.SwitchOn
		param.CommonTime = enum_agent.HeartbeatCommonTime
		param.FastTime = enum_agent.HeartbeatFastTime
		param.AutoWhiteSwitch = enums.SwitchOn
		param.AlarmAiSwitch = enums.SwitchOn
		param.LoginSolveSwitch = enums.SwitchOff
		param.LoginSolveConfig = jsonLoginSolve
		param.HardwareSwitch = enums.SwitchOn
		param.HardwareTime = enum_agent.HardwareTime
		param.CreatedAt = nowDate

		isOk := crd.masterModel.InsertItem(param)
		if !isOk {
			return false, "新建数据失败"
		}

		//添加任务
		crd.addTask(param)

		//自动授权
		crd.autoGrant(param.AgentId)

	} else {
		updates := map[string]interface{}{
			"host_name":      param.HostName,
			"mac_id":         param.MacId,
			"agent_ip":       param.AgentIp,
			"net_ip":         param.NetIp,
			"os_name":        param.OsName,
			"platform_type":  param.PlatformType,
			"os_type":        param.OsType,
			"mem_desc":       param.MemDesc,
			"hd_desc":        param.HdDesc,
			"cpu_desc":       param.CpuDesc,
			"network_card":   param.NetworkCard,
			"agent_start_at": param.AgentStartAt,
			"os_start_at":    param.OsStartAt,
			"ver":            param.Ver,
			"is_online":      enums.IsYes,
			"online_at":      nowDate,
			"live_at":        nowDate,
		}
		isOk := crd.masterModel.UpdateById(param.AgentId, updates)
		if !isOk {
			return false, "更新数据失败"
		}

		//添加任务
		crd.addTask(agentModel)

		//有组织架构
		if agentModel.DeptId > 0 {
			isAllocateDept = false
		}
	}

	//自动分配组织架构
	if isAllocateDept {
		crd.autoAllocateDept(param.AdminId, param.AgentId, param.AgentIp)
	}

	//生成token
	token := CreateCacheFactory().EncryptByLogin(param)

	return true, token
}

// autoGrant 自动授权
func (crd *AgentAuth) autoGrant(agentId string) {
	go grant_service.CreateCurdFactory().AutoGrant(agentId)
}

// autoAllocateDept 自动分配组织架构
func (crd *AgentAuth) autoAllocateDept(adminId int64, agentId, agentIp string) {
	go department_service.CreateCurdFactory().AutoAllocate(adminId, agentId, agentIp)
}
