package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/model_organization"
)

type Org struct {
	transformer.Common
}

func (t *Org) pageGeneral(model *model_organization.OrganizationModel) map[string]interface{} {
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
func (t *Org) TransPage(total int64, models []*model_organization.OrganizationModel) interface{} {
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

// TransDetail 转换成详情页需要的数据
func (t *Org) TransDetail(model *model_organization.OrganizationModel) map[string]interface{} {
	return t.pageGeneral(model)
}
