package agent_service

import (
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers"
)

type SearchConditions struct {
	AgentId      string
	HostName     string
	AgentIp      string
	NetIp        string
	PlatformType int
	IsOnline     int
	IsSafe       int
	Alias        string
	ContactName  string
	DeptIds      []int
}

type AgentBase struct {
}

// search2Where 查询条件转化
func (crd *AgentBase) search2Where(search SearchConditions) (wheres map[string]interface{}) {
	if !helpers.StringIsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if !helpers.StringIsEmpty(search.HostName) {
		wheres["host_name"] = search.HostName
	}
	if !helpers.StringIsEmpty(search.NetIp) {
		wheres["net_ip"] = search.NetIp
	}
	if !helpers.StringIsEmpty(search.AgentIp) {
		wheres["agent_ip"] = gorm.Expr("LIKE ?", "%"+search.AgentIp+"%")
	}
	if search.PlatformType > 0 {
		wheres["platform_type"] = search.PlatformType
	}
	if search.IsOnline != enums.SearchAll {
		wheres["is_online"] = search.IsOnline
	}
	if search.IsSafe != enums.SearchAll {
		wheres["is_safe"] = search.IsSafe
	}
	if !helpers.StringIsEmpty(search.ContactName) {
		wheres["contact_name"] = search.ContactName
	}
	if !helpers.StringIsEmpty(search.Alias) {
		wheres["alias"] = search.Alias
	}
	if search.DeptIds != nil && len(search.DeptIds) > 0 {
		wheres["dept_id"] = search.DeptIds
	}
	return wheres
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

// addTask 新增任务
func (crd *AgentBase) addTask(param *model_agent.AgentModel) {
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
