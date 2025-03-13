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
	InternetSwitch      *int
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
func (crd *AgentBase) getUpdates(model *agent_models.AgentModel, param *ManualColumn) (err error, hasTask, hasConfigTask bool,
	updates map[string]interface{}, openAutoWhite bool) {
	if model == nil {
		return errors.New("终端不存在"), hasTask, hasConfigTask, updates, openAutoWhite
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

	if param.SelfProtectSwitch != nil && *param.SelfProtectSwitch != enums.AllDefaultInt {
		updates["self_protect_switch"] = *param.SelfProtectSwitch
		hasTask = true
		hasConfigTask = true
	}
	if param.InternetSwitch != nil && *param.InternetSwitch != enums.AllDefaultInt {
		updates["internet_switch"] = *param.InternetSwitch
		hasTask = true
		hasConfigTask = true
	}

	//如果未授权, 下面功能不可用
	if model.GrantStatus == enum_agent.GrantStatusUnauthorized {
		return nil, hasTask, hasConfigTask, updates, openAutoWhite
	}

	//以下授权可用功能
	grantConfig := enum_product.Json2Config(model.GrantConfig)
	if param.CommonTime > 0 && model.CommonTime != param.CommonTime {
		updates["common_time"] = param.CommonTime
		hasTask = true
		hasConfigTask = true
	}
	if param.FastTime > 0 && model.FastTime != param.FastTime {
		updates["fast_time"] = param.FastTime
		hasTask = true
		hasConfigTask = true
	}

	//可信防护已授权
	if grantConfig.Process == enums.SwitchOn {
		//自动白名单
		if param.AutoWhiteSwitch != nil && *param.AutoWhiteSwitch != enums.AllDefaultInt {
			updates["auto_white_switch"] = *param.AutoWhiteSwitch
			//开启自动匹配白名单
			if model.AutoWhiteSwitch != *param.AutoWhiteSwitch && *param.AutoWhiteSwitch == enums.SwitchOn {
				openAutoWhite = true
			}
		}

		//进程相关
		if param.ProcessSwitch != nil && *param.ProcessSwitch != enums.AllDefaultInt {
			updates["process_switch"] = *param.ProcessSwitch
			//有任务
			if model.ProcessSwitch != *param.ProcessSwitch {
				hasTask = true
				hasConfigTask = true
			}
			//方便下发相关任务
			model.ProcessSwitch = *param.ProcessSwitch
			//开关关闭模式默认成监控
			if *param.ProcessSwitch == enums.SwitchOff {
				updates["process_mode"] = enum_agent.ModeMonitor
			}
		}
		if param.ProcessMode > 0 {
			if model.ProcessSwitch == enums.SwitchOn {
				updates["process_mode"] = param.ProcessMode
				if model.ProcessMode != param.ProcessMode {
					hasTask = true
					hasConfigTask = true
				}
			} else {
				updates["process_mode"] = enum_agent.ModeMonitor
			}
		}

		//模块相关
		if model.PlatformType == enum_software.PlatformTypeWindows {
			if param.DllSwitch != nil && *param.DllSwitch != enums.AllDefaultInt {
				updates["dll_switch"] = *param.DllSwitch
				//有任务
				if model.DllSwitch != *param.DllSwitch {
					hasTask = true
					hasConfigTask = true
				}
				//方便下发相关任务
				model.DllSwitch = *param.DllSwitch
				//开关关闭模式默认成监控
				if *param.DllSwitch == enums.SwitchOff {
					updates["dll_mode"] = enum_agent.ModeMonitor
				}
			}
			if param.DllMode > 0 {
				if model.DllSwitch == enums.SwitchOn {
					updates["dll_mode"] = param.DllMode
					if model.DllMode != param.DllMode {
						hasTask = true
						hasConfigTask = true
					}
				} else {
					updates["dll_mode"] = enum_agent.ModeMonitor
				}
			}
		}

	}

	//防篡改已授权
	if grantConfig.Tamper == enums.SwitchOn {
		if param.TamperSwitch != nil && *param.TamperSwitch != enums.AllDefaultInt {
			updates["tamper_switch"] = *param.TamperSwitch
			//有任务
			if model.TamperSwitch != *param.TamperSwitch {
				hasTask = true
				hasConfigTask = true
			}
			//方便下发相关任务
			model.TamperSwitch = *param.TamperSwitch
			//开关关闭模式默认成监控
			if *param.TamperSwitch == enums.SwitchOff {
				updates["tamper_mode"] = enum_agent.ModeMonitor
			}
		}
		if param.TamperMode > 0 {
			if model.TamperSwitch == enums.SwitchOn {
				updates["tamper_mode"] = param.TamperMode
				if model.TamperMode != param.TamperMode {
					hasTask = true
					hasConfigTask = true
				}
			} else {
				updates["tamper_mode"] = enum_agent.ModeMonitor
			}
		}
	}

	//勒索防护已授权
	if grantConfig.Blackmail == enums.SwitchOn {
		if param.BlackmailSwitch != nil && *param.BlackmailSwitch != enums.AllDefaultInt {
			updates["blackmail_switch"] = *param.BlackmailSwitch
			//有任务
			if model.BlackmailSwitch != *param.BlackmailSwitch {
				hasTask = true
				hasConfigTask = true
			}
			//方便下发相关任务
			model.BlackmailSwitch = *param.BlackmailSwitch
			//开关关闭模式默认成监控
			if *param.BlackmailSwitch == enums.SwitchOff {
				updates["blackmail_mode"] = enum_agent.ModeMonitor
			}
		}
		if param.BlackmailMode > 0 {
			if model.BlackmailSwitch == enums.SwitchOn {
				updates["blackmail_mode"] = param.BlackmailMode
				if model.BlackmailMode != param.BlackmailMode {
					hasTask = true
					hasConfigTask = true
				}
			} else {
				updates["blackmail_mode"] = enum_agent.ModeMonitor
			}
		}
	}

	//端口虚开已授权
	if model.OsType != enum_software.OsTypeWinXp {
		if grantConfig.Port == enums.SwitchOn {
			if param.PortSwitch != nil && *param.PortSwitch != enums.AllDefaultInt {
				updates["port_switch"] = *param.PortSwitch
				//有任务
				if model.PortSwitch != *param.PortSwitch {
					hasTask = true
					hasConfigTask = true
				}
				//方便下发相关任务
				model.PortSwitch = *param.PortSwitch
			}
		}
	}

	//封堵已授权
	if model.OsType != enum_software.OsTypeWinXp {
		if grantConfig.Plug == enums.SwitchOn {
			if param.PlugSwitch != nil && *param.PlugSwitch != enums.AllDefaultInt {
				updates["plug_switch"] = *param.PlugSwitch
				//有任务
				if model.PlugSwitch != *param.PlugSwitch {
					hasTask = true
					hasConfigTask = true
				}
				//方便下发相关任务
				model.PlugSwitch = *param.PlugSwitch
			}
		}
	}

	//外设已授权
	if grantConfig.Usb == enums.SwitchOn {
		if param.UsbSwitch != nil && *param.UsbSwitch != enums.AllDefaultInt {
			updates["usb_switch"] = *param.UsbSwitch
			//有任务
			if model.UsbSwitch != *param.UsbSwitch {
				hasTask = true
				hasConfigTask = true
			}
			//方便下发相关任务
			model.UsbSwitch = *param.UsbSwitch
			//开关关闭模式默认成监控
			if *param.UsbSwitch == enums.SwitchOff {
				updates["usb_mode"] = enum_agent.ModeMonitor
			}
		}
		if param.UsbMode > 0 {
			if model.UsbSwitch == enums.SwitchOn {
				updates["usb_mode"] = param.UsbMode
				if model.UsbMode != param.UsbMode {
					hasTask = true
					hasConfigTask = true
				}
			} else {
				updates["usb_mode"] = enum_agent.ModeMonitor
			}
		}
	}

	//通讯日志已授权
	if grantConfig.Syslog == enums.SwitchOn {
		if param.SyslogDnsSwitch != nil && *param.SyslogDnsSwitch != enums.AllDefaultInt {
			updates["syslog_dns_switch"] = *param.SyslogDnsSwitch
			//有任务
			if model.SyslogDnsSwitch != *param.SyslogDnsSwitch {
				hasTask = true
				hasConfigTask = true
			}
		}
		if param.SyslogOuterSwitch != nil && *param.SyslogOuterSwitch != enums.AllDefaultInt {
			updates["syslog_outer_switch"] = *param.SyslogOuterSwitch
			//有任务
			if model.SyslogOuterSwitch != *param.SyslogOuterSwitch {
				hasTask = true
				hasConfigTask = true
			}
		}
		if param.SyslogInnerSwitch != nil && *param.SyslogInnerSwitch != enums.AllDefaultInt {
			updates["syslog_inner_switch"] = *param.SyslogInnerSwitch
			//有任务
			if model.SyslogInnerSwitch != *param.SyslogInnerSwitch {
				hasTask = true
				hasConfigTask = true
			}
		}
		if param.SyslogProcessSwitch != nil && *param.SyslogProcessSwitch != enums.AllDefaultInt {
			updates["syslog_process_switch"] = *param.SyslogProcessSwitch
			//有任务
			if model.SyslogProcessSwitch != *param.SyslogProcessSwitch {
				hasTask = true
				hasConfigTask = true
			}
		}
		if param.SyslogLoginSwitch != nil && *param.SyslogLoginSwitch != enums.AllDefaultInt {
			updates["syslog_login_switch"] = *param.SyslogLoginSwitch
			//有任务
			if model.SyslogLoginSwitch != *param.SyslogLoginSwitch {
				hasTask = true
				hasConfigTask = true
			}
		}
	}

	//外联已授权
	if grantConfig.Outreach == enums.SwitchOn {
		if param.OutreachAlarmSwitch != nil && *param.OutreachAlarmSwitch != enums.AllDefaultInt {
			updates["outreach_alarm_switch"] = *param.OutreachAlarmSwitch
			//有任务
			if model.OutreachAlarmSwitch != *param.OutreachAlarmSwitch {
				hasTask = true
				hasConfigTask = true
			}
		}
	}

	//口令破解已授权
	if grantConfig.Login == enums.SwitchOn {
		if param.LoginSolveSwitch != nil && *param.LoginSolveSwitch != enums.AllDefaultInt {
			updates["login_solve_switch"] = *param.LoginSolveSwitch
		}
	}

	return err, hasTask, hasConfigTask, updates, openAutoWhite
}

// addTask 新增任务
func (crd *AgentBase) addTask(model *agent_models.AgentModel, param *ManualColumn, hasConfigTask, isTrustDir bool, remark string) {
	var tasks []int
	if hasConfigTask {
		tasks = append(tasks, enum_task.TypeAgentConfigDownload)
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

	/*//上传准入开启
	if param.InternetSwitch != nil && *param.InternetSwitch == enums.SwitchOn {
		//windows XP和2003  不支持
		if model.OsType == enum_software.OsTypeWinXp {
			tasks = append(tasks, enum_task.TypePluggingDownload, enum_task.TypePlugIpBlackDownload)
		}
	}*/

	//外设开启
	if param.UsbSwitch != nil && *param.UsbSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeUsbWhiteDownload, enum_task.TypeUsbBlackDownload)
	}

	//外联管控开启
	if param.OutreachAlarmSwitch != nil && *param.OutreachAlarmSwitch == enums.SwitchOn {
		tasks = append(tasks, enum_task.TypeOutreachDetectDownload)
	}

	task_service.CreateCurdFactory().AddTasksForAgent(model.OrgId, model.AgentId, tasks, remark)

}
