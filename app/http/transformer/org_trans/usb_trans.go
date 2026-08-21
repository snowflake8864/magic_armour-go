package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_usb"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/usb_models"
)

type Usb struct {
	transformer.Common
	Lang string
}

// TransGlobalPage 转换成全局名单分页需要的数据
func (t *Usb) TransGlobalPage(total int64, models []*usb_models.GlobalModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["name"] = model.Name
			data["intro"] = model.Intro
			data["type"] = model.Type
			data["status"] = model.Status
			data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_usb.EnumMapInt)
			data["alias"] = model.Alias
			data["remark"] = model.Remark
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
func (t *Usb) TransPolicyPage(total int64, models []*usb_models.UsbModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["name"] = model.Name
			data["intro"] = model.Intro
			data["type"] = model.Type
			data["is_online"] = model.IsOnline
			data["status"] = model.Status
			data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_usb.EnumMapInt)
			data["alias"] = model.Alias
			data["remark"] = model.Remark
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

// TransLogPage 转换成日志页需要的数据
func (t *Usb) TransLogPage(total int64, models []*usb_models.LogModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["usb_hash"] = model.UsbHash
			data["name"] = model.Name

			data["status"] = model.Status
			data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_usb.EnumMapInt)
			data["plug_status"] = model.PlugStatus
			data["plug_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_status", model.PlugStatus, enum_usb.EnumMapInt)
			data["created_at"] = model.CreatedAt

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
