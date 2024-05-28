package agent_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateOnlineFactory() *AgentOnline {
	return &AgentOnline{masterModel: agent_models.CreateAgentFactory()}
}

type AgentOnline struct {
	AgentBase
	masterModel *agent_models.AgentModel
}

// RegisterOrLogin 存在登录，不存在注册
func (crd *AgentOnline) RegisterOrLogin(param *agent_models.AgentModel) (bool, string) {
	agentModel := crd.masterModel.GetItemByOrgAndAgent(param.OrgId, param.AgentId, crd.masterModel.ColumnSwitchAndModeForAlarm(true))
	nowDate := date_helper.CurDatetime()
	isAllocateDept := true //是否分配组织架构

	if string_helper.IsEmpty(agentModel.AgentId) {
		//授权配置
		grantConfig := enum_product.ProductGrantConfig{}
		grantConfig.Default()

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
		param.GrantConfig = json_helper.Encode(grantConfig)
		param.GrantStatus = enums.IsNot
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
		param.HardwareSwitch = enums.SwitchOn
		param.HardwareTime = enum_agent.HardwareTime
		param.CreatedAt = nowDate

		isOk := crd.masterModel.InsertItem(param)
		if !isOk {
			return false, "新建数据失败"
		}

		//添加扩展数据
		CreateExtendFactory().AddDefault(param.OrgId, param.AgentId)

		//添加任务
		crd.addTask(param, true, "终端验证不存在")

		//自动授权
		crd.autoGrant(param.OrgId, param.AgentId)

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
		isOk, _ := crd.masterModel.UpdateByOrgAndAgent(param.OrgId, param.AgentId, updates)
		if !isOk {
			return false, "更新数据失败"
		}

		//添加任务
		crd.addTask(agentModel, false, "终端验证已存在")

		//有组织架构
		if agentModel.DeptId > 0 {
			isAllocateDept = false
		}
	}

	//自动分配组织架构,使用通信IP
	if isAllocateDept {
		crd.autoAllocateDept(param.OrgId, param.AgentId, param.NetIp)
	}

	//生成token
	token := CreateCacheFactory().EncryptByLogin(param)

	return true, token
}

// autoGrant 自动授权
func (crd *AgentOnline) autoGrant(orgId int64, agentId string) {

}

// autoAllocateDept 自动分配组织架构
func (crd *AgentOnline) autoAllocateDept(orgId int64, agentId, agentIp string) {
}
