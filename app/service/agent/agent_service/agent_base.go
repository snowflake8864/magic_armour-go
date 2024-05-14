package agent_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_model"
	"goskeleton/app/service/agent/task_service"
)

type UpdateColumn struct {
	SyslogDnsSwitch     int `json:"syslog_dns_switch"`
	SyslogOuterSwitch   int `json:"syslog_outer_switch"`
	SyslogInnerSwitch   int `json:"syslog_inner_switch"`
	SyslogProcessSwitch int `json:"syslog_process_switch"`
	SyslogLoginSwitch   int `json:"syslog_login_switch"`
	OutreachAlarmSwitch int `json:"outreach_alarm_switch"`
	SelfProtectSwitch   int `json:"self_protect_switch"`
	AutoWhiteSwitch     int `json:"auto_white_switch"`
}

type AgentBase struct {
}

// addTask 新增任务
func (crd *AgentBase) addTask(param *agent_model.AgentModel, isTrustDir bool, remark string) {
	tasks := []int{
		enum_task.TypeAgentConfigDownload,
	}
	if isTrustDir {
		tasks = append(tasks, enum_task.TypeTrustDir)
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

	//外设开启
	if param.UsbSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeUsbWhiteDownload, enum_task.TypeUsbBlackDownload)
	}

	task_service.CreateCurdFactory().AddTasksForAgent(param.OrgId, param.AgentId, tasks, remark)

}
