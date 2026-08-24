package org_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_port"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/helpers/struct_helper"
)

type PortTrans struct {
	transformer.Common
	Lang string
}

// TransGlobalPage 转换成全局名单分页需要的数据
func (t *PortTrans) TransGlobalPage(total int64, models []*port_models.VirtualGlobalModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["dept_fullname"] = t.GetDeptFullname(model.DeptId, departments)

			if t.Lang == consts.LanguageZh {
				data["protocol_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("protocol", model.Protocol, enum_port.EnumMapStr)
				data["type_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("ip_type", model.Type, enum_port.EnumMapStr)
				data["forward_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("forward_type", model.ForwardType, enum_port.EnumMapInt)
				data["alarm_level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("alarm_level", model.AlarmLevel, enum_port.EnumMapInt)
			} else {
				data["protocol_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("protocol_en", model.Protocol, enum_port.EnumMapStr)
				data["type_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("ip_type_en", model.Type, enum_port.EnumMapStr)
				data["forward_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("forward_type_en", model.ForwardType, enum_port.EnumMapInt)
				data["alarm_level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("alarm_level_en", model.AlarmLevel, enum_port.EnumMapInt)
			}

			forwardPort := ""
			if string_helper.ToInt(model.ForwardPort) > 0 {
				forwardPort = model.ForwardPort
			}
			data["forward_port"] = forwardPort

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
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["protocol_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("protocol", model.Protocol, enum_port.EnumMapStr)
				data["type_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("ip_type", model.Type, enum_port.EnumMapStr)
				data["forward_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("forward_type", model.ForwardType, enum_port.EnumMapInt)
				data["alarm_level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("alarm_level", model.AlarmLevel, enum_port.EnumMapInt)
			} else {
				data["protocol_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("protocol_en", model.Protocol, enum_port.EnumMapStr)
				data["type_label"] = enums.Get1stIsStrAnd2ndIsStrReturnStr("ip_type_en", model.Type, enum_port.EnumMapStr)
				data["forward_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("forward_type_en", model.ForwardType, enum_port.EnumMapInt)
				data["alarm_level_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("alarm_level_en", model.AlarmLevel, enum_port.EnumMapInt)
			}

			forwardPort := ""
			if string_helper.ToInt(model.ForwardPort) > 0 {
				forwardPort = model.ForwardPort
			}
			data["forward_port"] = forwardPort

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
			data["id"] = fmt.Sprintf("%d", model.Id)
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
