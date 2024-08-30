package agent_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/task_service"
)

// ManualColumn 人工可干预字段
type ManualColumn struct {
	UserId              *int64
	DeptId              *int64
	AgentType           *int
	Alias               *string
	ContactName         *string
	ContactTel          *string
	Position            *string
	Remark              *string
	ProcessSwitch       *int
	ProcessMode         int
	DllSwitch           *int
	DllMode             int
	TamperSwitch        *int
	TamperMode          int
	BlackmailSwitch     *int
	BlackmailMode       int
	PortSwitch          *int
	PlugSwitch          *int
	UsbSwitch           *int
	UsbMode             int
	SyslogDnsSwitch     *int
	SyslogOuterSwitch   *int
	SyslogInnerSwitch   *int
	SyslogProcessSwitch *int
	SyslogLoginSwitch   *int
	OutreachAlarmSwitch *int
	SelfProtectSwitch   *int
	KernelUpgradeSwitch *int
	CommonTime          int
	FastTime            int
	AutoWhiteSwitch     *int
	AlarmAiSwitch       *int
	LoginSolveSwitch    *int
	HardwareSwitch      *int
	HardwareTime        int
}

type AgentBase struct {
}

// 获取修改数据
func (crd *AgentBase) getUpdates(model *agent_models.AgentModel, param *ManualColumn) (err error, hasTask bool,
	openAutoWhite bool, updates map[string]interface{}) {
	if model == nil {
		return errors.New("终端不存在"), hasTask, openAutoWhite, updates
	}
	updates = make(map[string]interface{})
	if param.DeptId != nil {
		updates["dept_id"] = 0
		if *param.DeptId != int64(enums.AllDefaultInt) {
			updates["dept_id"] = *param.DeptId
		}
	}
	if param.AgentType != nil {
		updates["agent_type"] = 0
		if *param.AgentType != enums.AllDefaultInt {
			updates["agent_type"] = *param.AgentType
		}
	}
	if param.Alias != nil && *param.Alias != enums.AllDefaultStr {
		updates["alias"] = *param.Alias
	}
	if param.ContactName != nil && *param.ContactName != enums.AllDefaultStr {
		updates["contact_name"] = *param.ContactName
	}
	if param.ContactTel != nil && *param.ContactTel != enums.AllDefaultStr {
		updates["contact_tel"] = *param.ContactTel
	}
	if param.Position != nil && *param.Position != enums.AllDefaultStr {
		updates["position"] = *param.Position
	}
	if param.Remark != nil && *param.Remark != enums.AllDefaultStr {
		updates["remark"] = *param.Remark
	}

	if param.SelfProtectSwitch != nil && *param.SelfProtectSwitch != enums.AllDefaultInt && model.SelfProtectSwitch != *param.SelfProtectSwitch {
		updates["self_protect_switch"] = *param.SelfProtectSwitch
		hasTask = true
	}
	//如果未授权, 下面功能不可用
	if model.GrantStatus == enum_agent.GrantStatusNot {
		return nil, hasTask, openAutoWhite, updates
	}

	//以下授权可用功能
	grantConfig := enum_product.Json2Config(model.GrantConfig)
	if param.CommonTime > 0 && model.CommonTime != param.CommonTime {
		updates["common_time"] = param.CommonTime
		hasTask = true
	}
	if param.FastTime > 0 && model.FastTime != param.FastTime {
		updates["fast_time"] = param.FastTime
		hasTask = true
	}

	//可信防护已授权
	if grantConfig.Process == enums.SwitchOn {
		//自动白名单
		if param.AutoWhiteSwitch != nil && *param.AutoWhiteSwitch != enums.AllDefaultInt && model.AutoWhiteSwitch != *param.AutoWhiteSwitch {
			updates["auto_white_switch"] = *param.AutoWhiteSwitch
			openAutoWhite = true
		}

		if param.ProcessSwitch != nil && *param.ProcessSwitch != enums.AllDefaultInt && model.ProcessSwitch != *param.ProcessSwitch {
			updates["process_switch"] = *param.ProcessSwitch
			hasTask = true
			//方便下发相关任务
			model.ProcessSwitch = *param.ProcessSwitch
			//同步模式
			if *param.ProcessSwitch == enums.SwitchOff {
				updates["process_mode"] = enum_agent.ModeMonitor
			}
		}
		if param.ProcessMode > 0 && model.ProcessMode != param.ProcessMode && model.ProcessSwitch == enums.SwitchOn {
			updates["process_mode"] = param.ProcessMode
			hasTask = true
		}

		//模块相关
		if model.PlatformType == enum_software.PlatformTypeWindows {
			if param.DllSwitch != nil && *param.DllSwitch != enums.AllDefaultInt && model.DllSwitch != *param.DllSwitch {
				updates["dll_switch"] = *param.DllSwitch
				hasTask = true
				//方便下发相关任务
				model.DllSwitch = *param.DllSwitch
				//同步模式
				if *param.DllSwitch == enums.SwitchOff {
					updates["dll_mode"] = enum_agent.ModeMonitor
				}
			}
			if param.DllMode > 0 && model.DllMode != param.DllMode && model.DllSwitch == enums.SwitchOn {
				updates["dll_mode"] = param.DllMode
				hasTask = true
			}
		}

	}

	//防篡改已授权
	if grantConfig.Tamper == enums.SwitchOn {
		if param.TamperSwitch != nil && *param.TamperSwitch != enums.AllDefaultInt && model.TamperSwitch != *param.TamperSwitch {
			updates["tamper_switch"] = *param.TamperSwitch
			hasTask = true
			//方便下发相关任务
			model.TamperSwitch = *param.TamperSwitch
			//同步模式
			if *param.TamperSwitch == enums.SwitchOff {
				updates["tamper_mode"] = enum_agent.ModeMonitor
			}
		}
		if param.TamperMode > 0 && model.TamperMode != param.TamperMode && model.TamperSwitch == enums.SwitchOn {
			updates["tamper_mode"] = param.TamperMode
			hasTask = true
		}
	}

	//勒索防护已授权
	if grantConfig.Blackmail == enums.SwitchOn {
		if param.BlackmailSwitch != nil && *param.BlackmailSwitch != enums.AllDefaultInt && model.BlackmailSwitch != *param.BlackmailSwitch {
			updates["blackmail_switch"] = *param.BlackmailSwitch
			hasTask = true
			//方便下发相关任务
			model.BlackmailSwitch = *param.BlackmailSwitch
			//同步模式
			if *param.BlackmailSwitch == enums.SwitchOff {
				updates["blackmail_mode"] = enum_agent.ModeMonitor
			}
		}
		if param.BlackmailMode > 0 && model.BlackmailMode != param.BlackmailMode && model.BlackmailSwitch == enums.SwitchOn {
			updates["blackmail_mode"] = param.BlackmailMode
			hasTask = true
		}
	}

	//端口虚开已授权
	if model.OsType != enum_software.OsTypeWinXp {
		if grantConfig.Port == enums.SwitchOn {
			if param.PortSwitch != nil && *param.PortSwitch != enums.AllDefaultInt && model.PortSwitch != *param.PortSwitch {
				updates["port_switch"] = *param.PortSwitch
				hasTask = true
				//方便下发相关任务
				if *param.PortSwitch == enums.SwitchOn {
					model.PortSwitch = enums.SwitchOn
				}
			}
		}
	}

	//封堵已授权
	if model.OsType != enum_software.OsTypeWinXp {
		if grantConfig.Plug == enums.SwitchOn {
			if param.PlugSwitch != nil && *param.PlugSwitch != enums.AllDefaultInt && model.PlugSwitch != *param.PlugSwitch {
				updates["plug_switch"] = *param.PlugSwitch
				hasTask = true
				//方便下发相关任务
				if *param.PlugSwitch == enums.SwitchOn {
					model.PlugSwitch = enums.SwitchOn
				}
			}
		}
	}

	//外设已授权
	if grantConfig.Usb == enums.SwitchOn {
		if param.UsbSwitch != nil && *param.UsbSwitch != enums.AllDefaultInt && model.UsbSwitch != *param.UsbSwitch {
			updates["usb_switch"] = *param.UsbSwitch
			hasTask = true
			//方便下发相关任务
			model.UsbSwitch = *param.UsbSwitch
			//同步模式
			if *param.UsbSwitch == enums.SwitchOff {
				updates["usb_mode"] = enum_agent.ModeMonitor
			}
		}
		if param.UsbMode > 0 && model.UsbMode != param.UsbMode && model.UsbSwitch == enums.SwitchOn {
			updates["usb_mode"] = param.UsbMode
			hasTask = true
		}
	}
	//通讯日志已授权
	if grantConfig.Syslog == enums.SwitchOn {
		if param.SyslogDnsSwitch != nil && *param.SyslogDnsSwitch != enums.AllDefaultInt && model.SyslogDnsSwitch != *param.SyslogDnsSwitch {
			updates["syslog_dns_switch"] = *param.SyslogDnsSwitch
			hasTask = true
		}
		if param.SyslogOuterSwitch != nil && *param.SyslogOuterSwitch != enums.AllDefaultInt && model.SyslogOuterSwitch != *param.SyslogOuterSwitch {
			updates["syslog_outer_switch"] = *param.SyslogOuterSwitch
			hasTask = true
		}
		if param.SyslogInnerSwitch != nil && *param.SyslogInnerSwitch != enums.AllDefaultInt && model.SyslogInnerSwitch != *param.SyslogInnerSwitch {
			updates["syslog_inner_switch"] = *param.SyslogInnerSwitch
			hasTask = true
		}
		if param.SyslogProcessSwitch != nil && *param.SyslogProcessSwitch != enums.AllDefaultInt && model.SyslogProcessSwitch != *param.SyslogProcessSwitch {
			updates["syslog_process_switch"] = *param.SyslogProcessSwitch
			hasTask = true
		}
		if param.SyslogLoginSwitch != nil && *param.SyslogLoginSwitch != enums.AllDefaultInt && model.SyslogLoginSwitch != *param.SyslogLoginSwitch {
			updates["syslog_login_switch"] = *param.SyslogLoginSwitch
			hasTask = true
		}
	}
	//外联已授权
	if grantConfig.Outreach == enums.SwitchOn {
		if param.OutreachAlarmSwitch != nil && *param.OutreachAlarmSwitch != enums.AllDefaultInt && model.OutreachAlarmSwitch != *param.OutreachAlarmSwitch {
			updates["outreach_alarm_switch"] = *param.OutreachAlarmSwitch
		}
	}
	//口令破解已授权
	if grantConfig.Login == enums.SwitchOn {
		if param.LoginSolveSwitch != nil && *param.LoginSolveSwitch != enums.AllDefaultInt && model.LoginSolveSwitch != *param.LoginSolveSwitch {
			updates["login_solve_switch"] = *param.LoginSolveSwitch
		}
	}

	return err, hasTask, openAutoWhite, updates
}

// addTask 新增任务
func (crd *AgentBase) addTask(model *agent_models.AgentModel, param *ManualColumn, isTrustDir bool, remark string) {
	tasks := []int{
		enum_task.TypeAgentConfigDownload,
	}
	if isTrustDir {
		tasks = append(tasks, enum_task.TypeTrustDir)
	}

	if param == nil {
		param = &ManualColumn{}
	}

	//进程开启
	if param.ProcessSwitch != nil && *param.ProcessSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeProcessWhiteDownload, enum_task.TypeProcessBlackDownload)
	}

	//模块开启
	if param.DllSwitch != nil && *param.DllSwitch == enums.SwitchOn {
		//todo linux不支持模块 http://192.168.19.9/zentao/story-view-636.html
		if model.PlatformType == enum_software.PlatformTypeWindows {
			tasks = append(tasks, enum_task.TypeDllWhiteDownload, enum_task.TypeDllBlackDownload)
		}
	}
	//防篡改开启
	if param.TamperSwitch != nil && *param.TamperSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeTamperPolicyDownload)
	}

	//勒索开启
	if param.BlackmailSwitch != nil && *param.BlackmailSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeBlackmailPolicyDownload)
	}

	//端口虚开开启
	if param.PortSwitch != nil && *param.PortSwitch == enums.SwitchOn {
		//windows XP和2003  不支持
		if model.OsType != enum_software.OsTypeWinXp {
			tasks = append(tasks, enum_task.TypePortWorkingUpload, enum_task.TypePortVirtualPolicyDownload)
		}
	}

	//动态封堵开启
	if param.PlugSwitch != nil && *param.PlugSwitch == enums.SwitchOn {
		//windows XP和2003  不支持
		if model.OsType == enum_software.OsTypeWinXp {
			tasks = append(tasks, enum_task.TypePluggingDownload, enum_task.TypePlugIpBlackDownload)
		}
	}

	//外设开启
	if param.UsbSwitch != nil && *param.UsbSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeUsbWhiteDownload, enum_task.TypeUsbBlackDownload)
	}

	task_service.CreateCurdFactory().AddTasksForAgent(model.OrgId, model.AgentId, tasks, remark)

}
