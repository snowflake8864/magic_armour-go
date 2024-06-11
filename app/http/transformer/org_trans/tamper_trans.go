package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/json_helper"
)

type TamperTrans struct {
	transformer.Common
}

// TransGlobalPage 转换成全局名单分页需要的数据
func (t *TamperTrans) TransGlobalPage(total int64, models []*tamper_models.GlobalModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["dir"] = model.Dir
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_tamper.EnumMapInt)
			data["platform_type"] = model.PlatformType
			data["platform_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("platform_type", model.PlatformType, enum_software.EnumMap)
			data["protect_rw"] = model.ProtectRw
			data["protect_rw_label"] = enum_tamper.GetPowerLabels(model.ProtectRw)
			data["include_suffix"] = model.IncludeSuffix
			data["exclude_suffix"] = model.ExcludeSuffix
			data["exclude_folder"] = model.ExcludeFolder
			data["is_extend"] = model.IsExtend
			data["is_extend_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("is_extend", model.IsExtend, enum_tamper.EnumMapInt)
			var processes []tamper_models.Process
			_ = json_helper.ToAny(model.Processes, &processes)
			data["processes"] = processes
			data["dept_id"] = model.DeptId
			data["dept_fullname"] = t.GetDeptFullname(model.DeptId, departments)

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
func (t *TamperTrans) TransPolicyPage(total int64, models []*tamper_models.PolicyModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["dir"] = model.Dir
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_tamper.EnumMapInt)
			data["protect_rw"] = model.ProtectRw
			data["protect_rw_label"] = enum_tamper.GetPowerLabels(model.ProtectRw)
			data["include_suffix"] = model.IncludeSuffix
			data["exclude_suffix"] = model.ExcludeSuffix
			data["exclude_folder"] = model.ExcludeFolder
			data["is_extend"] = model.IsExtend
			data["is_extend_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("is_extend", model.IsExtend, enum_tamper.EnumMapInt)
			var processes []tamper_models.Process
			_ = json_helper.ToAny(model.Processes, &processes)
			data["processes"] = processes

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
