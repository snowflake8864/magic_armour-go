package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/organization/organization_models"
)

type Org struct {
	transformer.Common
}

func (t *Org) pageGeneral(model *organization_models.OrganizationModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["code"] = model.Code
	data["name"] = model.Name
	data["type"] = model.Type
	data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_organization.EnumMap)
	data["user_name"] = model.UserName
	data["user_email"] = model.UserEmail
	data["user_phone"] = model.UserPhone
	data["status"] = model.Status
	data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_organization.EnumMap)
	data["is_clean"] = model.IsClean
	data["is_del"] = model.IsDel
	data["huanshi_switch"] = model.HuanshiSwitch
	data["trusteeship_switch"] = model.TrusteeshipSwitch
	data["login_phone_check"] = model.LoginPhoneCheck
	data["remark"] = model.Remark
	data["source"] = model.Source
	data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source", model.Source, enums.EnumMapInt)
	data["created_at"] = model.CreatedAt
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
func (t *Org) TransOtherPage(total int64, models []*organization_models.OrganizationModel, extendInfo map[int64]interface{}) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)

			//扩展信息
			var einfo interface{}
			tmp["extend_info"] = einfo
			if extendInfo != nil && len(extendInfo) > 0 {
				if extend, ok := extendInfo[model.Id]; ok {
					tmp["extend_info"] = extend
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
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["module"] = model.Module
	data["module_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("log_action", model.Module, enum_organization.EnumMap)
	data["org_id"] = model.OrgId
	data["org_name"] = model.OrgName
	data["user_name"] = model.UserName
	data["user_email"] = model.UserEmail
	data["user_phone"] = model.UserPhone
	data["content"] = model.Content
	data["operate_source"] = model.OperateSource
	data["operate_source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source", model.OperateSource, enums.EnumMapInt)
	data["operate_id"] = model.OperateId
	data["operate_name"] = model.OperateName
	data["created_at"] = model.CreatedAt
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
