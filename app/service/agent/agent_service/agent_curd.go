package agent_service

import (
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_grant"
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/department_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/grant_service"
	"goskeleton/app/utils/helpers"
)

func CreateCurdFactory() *AgentCurd {
	return &AgentCurd{masterModel: model_agent.CreateAgentFactory()}
}

type AgentCurd struct {
	masterModel *model_agent.AgentModel
}

// ColumnConfig 字段:全部配置相关
func (crd *AgentCurd) ColumnConfig(needId bool) string {
	column := ""
	if needId {
		column = "agent_id,"
	}
	column = column +
		"process_switch, process_mode, dll_switch, dll_mode, tamper_switch, tamper_mode, " +
		"blackmail_switch, blackmail_mode, port_switch, plug_switch, usb_switch, usb_mode," +
		"syslog_dns_switch, syslog_outer_switch, syslog_inner_switch, syslog_process_switch, syslog_login_switch," +
		"outreach_alarm_switch, self_protect_switch, kernel_upgrade_switch, auto_white_switch, alarm_ai_switch," +
		"login_solve_switch, hardware_switch, hardware_time, common_time,fast_time"
	return column
}

// ColumnSwitchAndModeForAlarm 字段:告警需求的开关和模式
func (crd *AgentCurd) ColumnSwitchAndModeForAlarm(needId bool) string {
	column := ""
	if needId {
		column = "agent_id,"
	}
	column = column +
		"process_switch, process_mode, dll_switch, dll_mode, tamper_switch, tamper_mode, " +
		"blackmail_switch, blackmail_mode, port_switch,plug_switch, usb_switch, usb_mode"
	return column
}

// RegisterOrLogin 存在登录，不存在注册
func (crd *AgentCurd) RegisterOrLogin(param *model_agent.AgentModel) (bool, string) {
	column := "agent_id,dept_id,platform_type,os_type,process_switch,dll_switch,tamper_switch,blackmail_switch,port_switch,plug_switch,usb_switch"
	agentModel := crd.GetItemByAgent(param.AgentId, column)
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
func (crd *AgentCurd) autoGrant(agentId string) {
	go grant_service.CreateCurdFactory().AutoGrant(agentId)
}

// autoAllocateDept 自动分配组织架构
func (crd *AgentCurd) autoAllocateDept(adminId int, agentId, agentIp string) {
	go department_service.CreateCurdFactory().AutoAllocate(adminId, agentId, agentIp)
}

// addTask 新增任务
func (crd *AgentCurd) addTask(param *model_agent.AgentModel) {
	tasks := []int{
		enum_task.TypeAgentConfigDownload,
		enum_task.TypeTrustDir,
	}

	//进程开启
	if param.ProcessSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeProcessWhiteDownload, enum_task.TypeProcessBlackDownload)
	}
	//模块开启
	if param.DllSwitch == enums.SwitchOn {
		//todo linux不支持模块 http://192.168.19.9/zentao/story-view-636.html
		if param.PlatformType == enum_software.PlatformTypeWindows {
			tasks = append(tasks, enum_task.TypeDllWhiteDownload, enum_task.TypeDllBlackDownload)
		}
	}
	//防篡改开启
	if param.TamperSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeTamperPolicyDownload)
	}

	//勒索开启
	if param.BlackmailSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeBlackmailPolicyDownload)
	}

	//端口虚开开启
	if param.PortSwitch == enums.SwitchOn {
		//windows XP和2003  不支持
		if param.OsType == enum_software.OsTypeWinXp {
			tasks = append(tasks, enum_task.TypePortWorkingUpload, enum_task.TypePortVirtualPolicyDownload)
		}
	}

	//动态封堵开启
	if param.PlugSwitch == enums.SwitchOn {
		//windows XP和2003  不支持
		if param.OsType == enum_software.OsTypeWinXp {
			tasks = append(tasks, enum_task.TypePluggingDownload, enum_task.TypePlugIpBlackDownload)
		}
	}

	//外设开机
	if param.UsbSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeUsbWhiteDownload, enum_task.TypeUsbBlackDownload)
	}

	task_service.CreateCurdFactory().AddTasksForAgent(param.AgentId, tasks, "login")

}

// GetItemByAgent 通过终端id获取一条记录
func (crd *AgentCurd) GetItemByAgent(agentId, column string) (model *model_agent.AgentModel) {
	model, _ = crd.masterModel.GetItemById(agentId, column)
	return model
}

// GetMapsByAgents 通过终端ids获取多条记录, 有映射关系
func (crd *AgentCurd) GetMapsByAgents(agentIds []string, column string) map[string]model_agent.AgentModel {
	models, err := crd.masterModel.GetItemsByIds(agentIds, column)
	if err != nil {
		return nil
	}
	if models == nil {
		return nil
	}
	result := make(map[string]model_agent.AgentModel)
	for _, model := range models {
		result[model.AgentId] = model
	}
	return result
}

// SetLivingByAgent 设置终端活动中
func (crd *AgentCurd) SetLivingByAgent(agentId string) bool {
	nowDate := helpers.GetDatetime()

	updates := map[string]interface{}{
		"live_at":   nowDate,
		"online_at": gorm.Expr("IF(is_online=?, ?, online_at)", enums.IsNot, nowDate),
		"is_online": enums.IsYes,
	}
	return crd.masterModel.UpdateById(agentId, updates)
}
