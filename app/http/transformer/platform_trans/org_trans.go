package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type Org struct {
	transformer.Common
	Lang string
}

func (t *Org) pageGeneral(model *organization_models.OrganizationModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	if t.Lang == consts.LanguageZh {
		data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_organization.EnumMap)
		data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_organization.EnumMap)
		data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source", model.Source, enums.EnumMapInt)
	} else {
		data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_organization.EnumMap)
		data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status_en", model.Status, enum_organization.EnumMap)
		data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source_en", model.Source, enums.EnumMapInt)
	}
	return data
}

// TransPage 转换成管理页需要的数据
func (t *Org) TransPage(total int64, models []*organization_models.OrganizationModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransOtherPage 转换成管理页需要的数据
func (t *Org) TransOtherPage(total int64, models []*organization_models.OrganizationModel, agentExtend map[int64]interface{}, orderExtend map[int64]interface{}) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)

			//扩展信息
			var einfo interface{}
			tmp["agent_info"] = einfo
			tmp["order_info"] = einfo
			if agentExtend != nil && len(agentExtend) > 0 {
				if extend, ok := agentExtend[model.Id]; ok {
					tmp["agent_info"] = extend
				}
			}
			if orderExtend != nil && len(orderExtend) > 0 {
				if extend, ok := orderExtend[model.Id]; ok {
					tmp["order_info"] = extend
				}
			}

			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransDetail 转换成详情页需要的数据
func (t *Org) TransDetail(model *organization_models.OrganizationModel) map[string]interface{} {
	return t.pageGeneral(model)
}

func (t *Org) pageLogGeneral(model *organization_models.LogModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)
	if t.Lang == consts.LanguageZh {
		data["module_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("log_action", model.Module, enum_organization.EnumMap)
		data["operate_source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source", model.OperateSource, enums.EnumMapInt)
	} else {
		data["module_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("log_action_en", model.Module, enum_organization.EnumMap)
		data["operate_source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source_en", model.OperateSource, enums.EnumMapInt)
	}
	return data
}

// TransLogPage 转换成日志页需要的数据
func (t *Org) TransLogPage(total int64, models []*organization_models.LogModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageLogGeneral(model)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransLogList 转换成日志列表
func (t *Org) TransLogList(models []*organization_models.LogModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageLogGeneral(model)
			items = append(items, tmp)
		}
	}

	return items
}
