package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

type Agent struct {
	transformer.Common
}

func (t *Agent) pageGeneral(model *agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel, extendInfo map[string]interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	data["agent_id"] = model.AgentId
	data["host_name"] = model.HostName
	data["mac_id"] = model.MacId
	data["agent_ip"] = model.AgentIp
	data["net_ip"] = model.NetIp
	data["os_name"] = model.OsName
	data["platform_type"] = model.PlatformType
	data["os_type"] = model.OsType
	data["is_online"] = model.IsOnline
	data["online_at"] = model.IsOnline
	data["is_safe"] = model.IsSafe
	data["contact_name"] = model.ContactName
	data["alias"] = model.Alias
	data["grant_config"] = json_helper.ToMap(model.GrantConfig)
	data["dept_id"] = model.DeptId
	data["dept_fullname"] = []interface{}{}
	if departments != nil && len(departments) > 0 {
		if deptModel, ok := departments[model.DeptId]; ok && !string_helper.IsEmpty(deptModel.FullName) {
			data["dept_fullname"] = json_helper.ToSliceString(deptModel.FullName)
		}
	}
	data["process_switch"] = model.ProcessSwitch
	data["process_mode"] = model.ProcessMode
	data["dll_switch"] = model.DllSwitch
	data["dll_mode"] = model.DllMode
	data["tamper_switch"] = model.TamperSwitch
	data["tamper_mode"] = model.TamperMode
	data["blackmail_switch"] = model.BlackmailSwitch
	data["blackmail_mode"] = model.BlackmailMode
	data["port_switch"] = model.PortSwitch
	data["plug_switch"] = model.PlugSwitch
	data["usb_switch"] = model.UsbSwitch
	data["usb_mode"] = model.UsbMode
	data["outreach_alarm_switch"] = model.OutreachAlarmSwitch
	data["login_solve_switch"] = model.LoginSolveSwitch
	data["created_at"] = model.CreatedAt

	//扩展信息
	var tmp interface{}
	data["extend_info"] = tmp
	if extendInfo != nil && len(extendInfo) > 0 {
		if extend, ok := extendInfo[model.AgentId]; ok {
			data["extend_info"] = extend
		}
	}

	return data
}

// TransManagePage 转换成管理分页需要的数据
func (t *Agent) TransManagePage(total int64, models []*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel, softwareNewest map[int]string) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model, departments, nil)

			canUpgrade := 0
			if isOk, _ := enum_software.CanUpgradeIndirect(softwareNewest, model.OsType, model.Ver); isOk {
				return 1
			}
			tmp["can_upgrade"] = canUpgrade
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransOtherPage 转换成其他分页需要的数据
func (t *Agent) TransOtherPage(total int64, models []*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel, extendInfo map[string]interface{}) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model, departments, extendInfo)

			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransList 转换列表需要的数据l
func (t *Agent) TransList(models []*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel, extendInfo map[string]interface{}) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model, departments, extendInfo)

			items = append(items, tmp)
		}
	}

	return items
}
