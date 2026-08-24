package org_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/struct_helper"
)

type TamperTrans struct {
	transformer.Common
	Lang string
}

// TransGlobalPage 转换成全局名单分页需要的数据
func (t *TamperTrans) TransGlobalPage(total int64, models []*tamper_models.GlobalModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["platform_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("platform_type", model.PlatformType, enum_software.EnumMap)
			data["dept_fullname"] = t.GetDeptFullname(model.DeptId, departments)

			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_tamper.EnumMapInt)
				data["protect_rw_label"] = enum_tamper.GetPowerLabels(model.ProtectRw)
				data["is_extend_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("is_extend", model.IsExtend, enum_tamper.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_tamper.EnumMapInt)
				data["protect_rw_label"] = enum_tamper.GetPowerEnLabels(model.ProtectRw)
				data["is_extend_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("is_extend_en", model.IsExtend, enum_tamper.EnumMapInt)
			}

			var processes []tamper_models.Process
			_ = json_helper.ToAny(model.Processes, &processes)
			data["processes"] = processes

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransPolicyPage 转换成策略页需要的数据
func (t *TamperTrans) TransPolicyPage(total int64, models []*tamper_models.PolicyModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := struct_helper.SimpleStructToMapWithJsonTags(model)
			data["id"] = fmt.Sprintf("%d", model.Id)

			if t.Lang == consts.LanguageZh {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_tamper.EnumMapInt)
				data["protect_rw_label"] = enum_tamper.GetPowerLabels(model.ProtectRw)
				data["is_extend_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("is_extend", model.IsExtend, enum_tamper.EnumMapInt)
			} else {
				data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.Type, enum_tamper.EnumMapInt)
				data["protect_rw_label"] = enum_tamper.GetPowerEnLabels(model.ProtectRw)
				data["is_extend_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("is_extend_en", model.IsExtend, enum_tamper.EnumMapInt)
			}

			var processes []tamper_models.Process
			_ = json_helper.ToAny(model.Processes, &processes)
			data["processes"] = processes

			items = append(items, data)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
