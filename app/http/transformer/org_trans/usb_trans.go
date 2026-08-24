package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_usb"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/utils/helpers/struct_helper"
)

type UsbTrans struct {
	transformer.Common
	Lang string
}

// TransGlobalPage 转换成全局名单分页需要的数据
func (t *UsbTrans) TransGlobalPage(total int64, models []*usb_models.GlobalModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)

			if t.Lang == consts.LanguageZh {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_usb.EnumMapInt)
			} else {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status_en", model.Status, enum_usb.EnumMapInt)
			}

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransPolicyPage 转换成策略页需要的数据
func (t *UsbTrans) TransPolicyPage(total int64, models []*usb_models.UsbModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)

			if t.Lang == consts.LanguageZh {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_usb.EnumMapInt)
			} else {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status_en", model.Status, enum_usb.EnumMapInt)
			}

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransLogPage 转换成日志页需要的数据
func (t *UsbTrans) TransLogPage(total int64, models []*usb_models.LogModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["usb_hash"] = model.UsbHash
			data["name"] = model.Name
			data["status"] = model.Status
			data["plug_status"] = model.PlugStatus
			data["created_at"] = model.CreatedAt

			if t.Lang == consts.LanguageZh {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_usb.EnumMapInt)
				data["plug_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_status", model.PlugStatus, enum_usb.EnumMapInt)
			} else {
				data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status_en", model.Status, enum_usb.EnumMapInt)
				data["plug_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("plug_status_en", model.PlugStatus, enum_usb.EnumMapInt)
			}
			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
