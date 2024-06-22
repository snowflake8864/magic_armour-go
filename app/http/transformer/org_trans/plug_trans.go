package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/model/organization/department_models"
)

type PlugTrans struct {
	transformer.Common
}

// TranIpPolicyPage 转换成终端ip名单分页需要的数据
func (t *PlugTrans) TranIpPolicyPage(total int64, models []*plug_models.IpPolicyModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["ip"] = model.Ip
			data["remark"] = model.Remark
			data["direction"] = model.Direction
			data["direction_label"] = enum_plug.GetDirectionLabel(model.Direction)
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

func (t *PlugTrans) dynamicGeneral(model *plug_models.DynamicModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["agent_id"] = model.AgentId
	data["plug_type"] = model.PlugType
	data["plug_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_type", model.PlugType, enum_plug.EnumMap)
	data["source"] = model.Source
	data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("source", model.Source, enum_plug.EnumMap)
	data["duration"] = model.Duration
	data["start_at"] = model.StartAt
	data["end_at"] = model.EndAt
	data["ip"] = model.Ip
	data["direction"] = model.Direction
	data["direction_label"] = enum_plug.GetDirectionLabel(model.Direction)
	data["open_switch"] = model.OpenSwitch
	data["created_at"] = model.CreatedAt
	data["update_id"] = model.UpdateId
	data["update_name"] = model.UpdateName
	data["updated_at"] = model.UpdatedAt

	return data
}

// TranDynamicPolicyPage 转换成终端动态封堵分页需要的数据
func (t *PlugTrans) TranDynamicPolicyPage(total int64, models []*plug_models.DynamicModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.dynamicGeneral(model)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TranDynamicGlobalPage 转换成全局动态封堵分页需要的数据
func (t *PlugTrans) TranDynamicGlobalPage(total int64, models []*plug_models.DynamicModel,
	agentModels map[string]*agent_models.AgentModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			data := t.dynamicGeneral(model)
			//终端信息

			data["agent_info"] = t.AppendAgentInfo(model.AgentId, agentModels, departments)
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TranIpGlobalPage 转换成全局ip名单分页需要的数据
func (t *PlugTrans) TranIpGlobalPage(total int64, models []*plug_models.IpGlobalModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["type"] = model.Type
			data["ip"] = model.Ip
			data["remark"] = model.Remark
			data["direction"] = model.Direction
			data["direction_label"] = enum_plug.GetDirectionLabel(model.Direction)
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
