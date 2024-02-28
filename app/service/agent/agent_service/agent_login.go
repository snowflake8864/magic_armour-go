package agent_service

import (
	"encoding/json"
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers"
)

func CreateLoginFactory() *AgentLogin {
	return &AgentLogin{
		masterModel: model_agent.CreateAgentFactory(),
	}
}

type AgentLogin struct {
	masterModel *model_agent.AgentModel
}

// RegisterOrLogin 存在登录，不存在注册
func (crd *AgentLogin) RegisterOrLogin(param *model_agent.AgentModel) (bool, string) {
	column := "agent_id,platform_type,os_type, process_switch, dll_switch,tamper_switch,blackmail_switch,port_switch,plug_switch,usb_switch"
	agentModel, _ := crd.masterModel.GetOneById(param.AgentId, column)
	fmt.Println("agentModel.AgentId=====================================")
	fmt.Println(agentModel.AgentId)
	nowDate := helpers.GetDatetime()
	isAllocateDept := true //是否分配组织架构

	if helpers.IsEmptyString(agentModel.AgentId) {
		plugPortConfig := enum_plug.PlugPortConfig{}
		plugPortConfig.Format()
		jsonPlugPort, _ := json.Marshal(plugPortConfig)

		plugLoginConfig := enum_plug.PlugLoginConfig{}
		plugLoginConfig.Format()
		jsonPlugLogin, _ := json.Marshal(plugLoginConfig)

		loginSolveConfig := enum_login_solve.LoginSolveConfig{}
		loginSolveConfig.Format()
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

		//添加任务
		crd.addTask(param)

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
			"cpu":            param.Cpu,
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
		crd.autoAllocateDept(param.AgentId, param.AdminId)
	}

	//生成token
	token := CreateCacheFactory().EncryptByLogin(param)

	return true, token
}

// autoGrant 自动授权
func (crd *AgentLogin) autoGrant(agentId string, adminId int) {

}

// autoAllocateDept 自动分配组织架构
func (crd *AgentLogin) autoAllocateDept(agentId string, adminId int) {

}

// addTask 新增任务
func (crd *AgentLogin) addTask(param *model_agent.AgentModel) {
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
