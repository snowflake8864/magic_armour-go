package agent_service

import (
	"encoding/json"
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
)

func CreateLoginFactory() *AgentLogin {
	return &AgentLogin{model_agent.CreateAgentFactory()}
}

type AgentLogin struct {
	masterModel *model_agent.AgentModel
}

// RegisterOrLogin 存在登录，不存在注册
func (crd *AgentLogin) RegisterOrLogin(param *model_agent.AgentModel) (bool, string) {
	column := "agent_id,os_type, process_switch, dll_switch,tamper_switch,blackmail_switch,port_switch,plug_switch,usb_switch"
	agentModel, _ := crd.masterModel.GetOneById(param.AgentId, column)
	fmt.Println(agentModel)
	nowDate := helpers.GetNowDate()
	isAllocateDept := true //是否分配组织架构

	if agentModel == nil {
		plugPortConfig := enum_plug.PlugPortConfig.Format
		jsonPlugPort, _ := json.Marshal(plugPortConfig)

		plugLoginConfig := enum_plug.PlugLoginConfig.Format
		jsonPlugLogin, _ := json.Marshal(plugLoginConfig)

		loginSolveConfig := enum_login_solve.LoginSolveConfig.Format
		jsonLoginSolve, _ := json.Marshal(loginSolveConfig)

		param.IsOnline = enums.IsYes
		param.OnlineAt = nowDate
		param.LiveAt = nowDate
		param.UpgradeVer = ""
		param.DeptId = 0
		param.AgentTypeId = 0
		param.Alias = ""
		param.ContactName = ""
		param.ContactTel = ""
		param.Remark = ""
		param.GrantConfig = ""
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
		param.PlugPortConfig = string(jsonPlugPort)
		param.PlugLoginConfig = string(jsonPlugLogin)
		param.UsbSwitch = enums.SwitchOff
		param.UsbMode = enum_agent.ModeMonitor
		param.SyslogDnsSwitch = enums.SwitchOff
		param.SyslogOuterSwitch = enums.SwitchOff
		param.SyslogInnerSwitch = enums.SwitchOff
		param.SyslogProcessSwitch = enums.SwitchOff
		param.SyslogLoginSwitch = enums.SwitchOff
		param.OutreachAlarmSwitch = enums.SwitchOff
		param.SelfProtectionSwitch = enums.SwitchOn
		param.KernelUpgradeSwitch = enums.SwitchOn
		param.CommonTime = enum_agent.HeartbeatCommonTime
		param.FastTime = enum_agent.HeartbeatFastTime
		param.AutoWhiteSwitch = enums.SwitchOff
		param.AlarmAiSwitch = enums.SwitchOff
		param.LoginSolveSwitch = enums.SwitchOff
		param.LoginSolveConfig = string(jsonLoginSolve)
		param.HardwareSwitch = enums.SwitchOn
		param.HardwareTime = enum_agent.HardwareTime
		param.CreatedAt = nowDate

		isOk := crd.masterModel.InsertItem(param)
		if !isOk {
			return false, "新建数据失败"
		}

		//自动授权
		crd.autoGrant(param.AgentId, param.AdminId)

	} else {
		updates := map[string]interface{}{
			"host_name":      param.HostName,
			"mac_id":         param.MacId,
			"agent_ip":       param.AgentIp,
			"net_ip":         param.NetIp,
			"os_name":        param.OsName,
			"platform_type":  param.PlatformType,
			"os_type":        param.OsType,
			"mem_size":       param.MemSize,
			"hd_size":        param.HdSize,
			"cpu_size":       param.CpuSize,
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

		//有组织架构
		if agentModel.DeptId > 0 {
			isAllocateDept = false
		}
	}

	//自动分配组织架构
	if isAllocateDept {

	}

	return true, ""
}

// autoGrant 自动授权
func (crd *AgentLogin) autoGrant(agentId string, adminId int) {

}

// autoAllocateDept 自动分配组织架构
func (crd *AgentLogin) autoAllocateDept(agentId string, adminId int) {

}
