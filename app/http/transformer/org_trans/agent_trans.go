package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/helpers/struct_helper"
)

type Agent struct {
	transformer.Common
}

func (t *Agent) pageGeneral(model *agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel, agentExtend map[string]interface{}) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	data["agent_id"] = model.AgentId
	data["host_name"] = model.HostName
	data["mac_id"] = model.MacId
	data["agent_ip"] = model.AgentIp
	data["net_ip"] = model.NetIp
	data["os_name"] = model.OsName
	data["platform_type"] = model.PlatformType
	data["os_type"] = model.OsType
	data["ver"] = model.Ver
	data["mem_desc"] = model.MemDesc
	data["hd_desc"] = model.HdDesc
	data["cpu_desc"] = model.CpuDesc
	data["network_card"] = model.CpuDesc
	data["agent_start_at"] = model.AgentStartAt
	data["os_start_at"] = model.OsStartAt

	data["is_online"] = model.IsOnline
	data["online_at"] = model.OnlineAt
	data["is_safe"] = model.IsSafe

	data["dept_id"] = model.DeptId
	data["dept_fullname"] = []interface{}{}
	if departments != nil && len(departments) > 0 {
		if deptModel, ok := departments[model.DeptId]; ok && !string_helper.IsEmpty(deptModel.FullName) {
			data["dept_fullname"] = json_helper.ToSliceString(deptModel.FullName)
		}
	}
	data["contact_name"] = model.ContactName
	data["contact_tel"] = model.ContactTel
	data["alias"] = model.Alias
	data["position"] = model.Position
	data["remark"] = model.Remark
	data["agent_type"] = model.AgentType

	productGrantConfig := enum_product.ProductGrantConfig{}
	_ = json_helper.ToStruct(model.GrantConfig, &productGrantConfig)
	data["grant_config"] = productGrantConfig
	data["grant_config_txt"] = productGrantConfig.SelectedTxt()
	data["grant_order_no"] = model.GrantOrderNo
	data["grant_product_title"] = model.GrantProductTitle
	data["grant_expire_at"] = model.GrantExpireAt
	data["grant_status"] = model.GrantStatus

	data["process_switch"] = model.ProcessSwitch
	data["process_mode"] = model.ProcessMode
	data["process_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode", model.ProcessMode, enum_agent.EnumIntMap)
	data["dll_switch"] = model.DllSwitch
	data["dll_mode"] = model.DllMode
	data["dll_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode", model.DllMode, enum_agent.EnumIntMap)
	data["tamper_switch"] = model.TamperSwitch
	data["tamper_mode"] = model.TamperMode
	data["tamper_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode", model.TamperMode, enum_agent.EnumIntMap)
	data["blackmail_switch"] = model.BlackmailSwitch
	data["blackmail_mode"] = model.BlackmailMode
	data["blackmail_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode", model.BlackmailMode, enum_agent.EnumIntMap)
	data["port_switch"] = model.PortSwitch
	data["plug_switch"] = model.PlugSwitch
	data["usb_switch"] = model.UsbSwitch
	data["usb_mode"] = model.UsbMode
	data["usb_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("mode", model.UsbMode, enum_agent.EnumIntMap)
	data["outreach_alarm_switch"] = model.OutreachAlarmSwitch
	data["login_solve_switch"] = model.LoginSolveSwitch
	data["syslog_dns_switch"] = model.SyslogDnsSwitch
	data["syslog_outer_switch"] = model.SyslogOuterSwitch
	data["syslog_inner_switch"] = model.SyslogInnerSwitch
	data["syslog_process_switch"] = model.SyslogProcessSwitch
	data["syslog_login_switch"] = model.SyslogLoginSwitch
	data["self_protect_switch"] = model.SelfProtectSwitch
	data["auto_white_switch"] = model.AutoWhiteSwitch
	data["common_time"] = model.CommonTime
	data["fast_time"] = model.FastTime
	data["kernel_upgrade_switch"] = model.KernelUpgradeSwitch
	data["hardware_switch"] = model.HardwareSwitch
	data["hardware_time"] = model.HardwareTime
	data["internet_switch"] = model.InternetSwitch
	data["created_at"] = model.CreatedAt

	//扩展信息
	var tmp interface{}
	data["extend_info"] = tmp
	if agentExtend != nil && len(agentExtend) > 0 {
		if extend, ok := agentExtend[model.AgentId]; ok {
			data["extend_info"] = extend
		}
	}

	return data
}

// TransManagePage 转换成管理分页需要的数据
func (t *Agent) TransManagePage(total int64, models []*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel, softwareNewest map[int]string) interface{} {
	var items []map[string]interface{}
	if models != nil {
		//终端类型
		typeJsonStr := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentTypes)
		agentTypes := enum_agent.GetTypeMap(typeJsonStr)

		for _, model := range models {
			tmp := t.pageGeneral(model, departments, nil)

			canUpgrade := 0
			if enum_agent.CanUpgrade(model.GrantStatus) {
				if isOk, _ := enum_software.CanUpgradeIndirect(softwareNewest, model.OsType, model.Ver); isOk {
					canUpgrade = 1
				}
			}
			tmp["can_upgrade"] = canUpgrade
			tmp["agent_type_label"] = enum_agent.GetTypeLabel(agentTypes, model.AgentType)
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

// TransSingle 转换单个需要的数据
func (t *Agent) TransSingle(model *agent_models.AgentModel) interface{} {
	return t.pageGeneral(model, nil, nil)
}
