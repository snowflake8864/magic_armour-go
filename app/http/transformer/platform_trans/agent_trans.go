package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/utils/helpers"
)

type Agent struct {
	Common
}

func (t *Agent) pageGeneral(model *model_agent.AgentModel, departments map[string]*model_platform.DepartmentModel) map[string]interface{} {
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
	data["grant_config"] = helpers.Json2Map(model.GrantConfig)
	data["dept_id"] = model.DeptId
	data["dept_fullname"] = []string{}
	if dept, ok := departments[model.AgentId]; ok && !helpers.StringIsEmpty(dept.FullName) {
		data["dept_fullname"] = helpers.Json2Map(model.GrantConfig)
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
	return data
}

// TransAgentManagePage 转换成管理页需要的数据
func (t *Agent) TransAgentManagePage(total int64, models []*model_agent.AgentModel, departments map[string]*model_platform.DepartmentModel, softwareNewest map[int]string) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model, departments)
			tmp["can_upgrade"] = 0
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
