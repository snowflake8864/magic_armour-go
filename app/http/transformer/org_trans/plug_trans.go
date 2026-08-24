package org_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type PlugTrans struct {
	transformer.Common
	Lang string
}

// TranIpPolicyPage 转换成终端ip名单分页需要的数据
func (t *PlugTrans) TranIpPolicyPage(total int64, models []*plug_models.IpPolicyModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["direction_label"] = enum_plug.GetDirectionLabel(model.Direction)

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

func (t *PlugTrans) dynamicGeneral(model *plug_models.DynamicModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	data["id"] = fmt.Sprintf("%d", model.Id)
	data["direction_label"] = enum_plug.GetDirectionLabel(model.Direction)

	if t.Lang == consts.LanguageZh {
		data["plug_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_type", model.PlugType, enum_plug.EnumMap)
		data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("source", model.Source, enum_plug.EnumMap)
	} else {
		data["plug_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_type", model.PlugType, enum_plug.EnumMap)
		data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("source", model.Source, enum_plug.EnumMap)
	}

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
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["direction_label"] = enum_plug.GetDirectionLabel(model.Direction)

			if t.Lang == consts.LanguageZh {
				data["effect_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("effect_mode", model.EffectMode, enum_plug.EnumMap)
			} else {
				data["effect_mode_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("effect_mode_en", model.EffectMode, enum_plug.EnumMap)
			}

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
