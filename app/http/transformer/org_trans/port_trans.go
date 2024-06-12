package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_port"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/model/organization/department_models"
)

type PortTrans struct {
	transformer.Common
}

// TransGlobalPage 转换成全局名单分页需要的数据
func (t *PortTrans) TransGlobalPage(total int64, models []*port_models.VirtualGlobalModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["protocol"] = model.Protocol
			data["protocol_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("protocol", model.Protocol, enum_port.EnumMapStr)
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("ip_type", model.Type, enum_port.EnumMapStr)
			data["dept_id"] = model.DeptId
			data["dept_fullname"] = t.GetDeptFullname(model.DeptId, departments)

			data["source_ip"] = model.SourceIp
			data["source_port"] = model.SourcePort
			data["forward_ip"] = model.ForwardIp
			data["forward_port"] = model.ForwardPort
			data["alarm_level"] = model.AlarmLevel
			data["alarm_level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("alarm_level", model.AlarmLevel, enum_port.EnumMapInt)
			data["open_switch"] = model.OpenSwitch
			data["create_id"] = model.CreateId
			data["create_name"] = model.CreateName
			data["created_at"] = model.CreatedAt
			data["update_id"] = model.UpdateId
			data["update_name"] = model.UpdateName
			data["updated_at"] = model.UpdatedAt

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransPolicyPage 转换成策略页需要的数据
func (t *PortTrans) TransPolicyPage(total int64, models []*port_models.VirtualPolicyModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["protocol"] = model.Protocol
			data["protocol_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("protocol", model.Protocol, enum_port.EnumMapStr)
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("ip_type", model.Type, enum_port.EnumMapStr)
			data["source_ip"] = model.SourceIp
			data["source_port"] = model.SourcePort
			data["forward_ip"] = model.ForwardIp
			data["forward_port"] = model.ForwardPort
			data["alarm_level"] = model.AlarmLevel
			data["alarm_level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("alarm_level", model.AlarmLevel, enum_port.EnumMapInt)
			data["open_switch"] = model.OpenSwitch

			data["create_id"] = model.CreateId
			data["create_name"] = model.CreateName
			data["created_at"] = model.CreatedAt
			data["update_id"] = model.UpdateId
			data["update_name"] = model.UpdateName
			data["updated_at"] = model.UpdatedAt

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransWorkingPage 转换成正常工作页需要的数据
func (t *PortTrans) TransWorkingPage(total int64, models []*port_models.WorkingModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["type"] = model.Type
			data["ip"] = model.Ip
			data["port"] = model.Port
			data["proc_id"] = model.ProcId
			data["proc_dir"] = model.ProcDir
			data["created_at"] = model.CreatedAt

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
