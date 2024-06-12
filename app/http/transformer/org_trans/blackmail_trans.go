package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_blackmail"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/blackmail_models"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/utils/helpers/json_helper"
)

type BlackmailTrans struct {
	transformer.Common
}

// TransGlobalPage 转换成全局名单分页需要的数据
func (t *BlackmailTrans) TransGlobalPage(total int64, models []*blackmail_models.GlobalModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_blackmail.EnumMapInt)
			data["type_val"] = model.TypeVal
			data["group_id"] = model.GroupId
			data["platform_type"] = model.PlatformType
			data["platform_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("platform_type", model.PlatformType, enum_software.EnumMap)
			var processes []tamper_models.Process
			_ = json_helper.ToAny(model.Processes, &processes)
			data["processes"] = processes
			data["dept_id"] = model.DeptId
			data["dept_fullname"] = t.GetDeptFullname(model.DeptId, departments)

			data["open_switch"] = model.OpenSwitch
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
func (t *BlackmailTrans) TransPolicyPage(total int64, models []*blackmail_models.PolicyModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["agent_id"] = model.AgentId
			data["type"] = model.Type
			data["type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.Type, enum_blackmail.EnumMapInt)
			data["type_val"] = model.TypeVal
			data["group_id"] = model.GroupId
			var processes []tamper_models.Process
			_ = json_helper.ToAny(model.Processes, &processes)
			data["processes"] = processes

			data["open_switch"] = model.OpenSwitch
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

// TransGroupPage 转换成文件类型组页需要的数据
func (t *BlackmailTrans) TransGroupPage(total int64, models []*blackmail_models.GroupModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["title"] = model.Title
			data["suffix"] = model.Suffix
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

// TransGroupList 转换成文件类型组页需要的数据
func (t *BlackmailTrans) TransGroupList(models []*blackmail_models.GroupModel) (items []map[string]interface{}) {
	if models != nil {
		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = model.Id
			data["suffix"] = model.Suffix
			data["title"] = model.Title
			/*data["create_id"] = model.CreateId
			data["create_name"] = model.CreateName
			data["created_at"] = model.CreatedAt
			data["update_id"] = model.UpdateId
			data["update_name"] = model.UpdateName
			data["updated_at"] = model.UpdatedAt*/

			items = append(items, data)
		}
	}
	return items
}
