package org_trans

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/agent/trust_dir_models"
	"goskeleton/app/model/organization/department_models"
)

type TrustDirTrans struct {
	transformer.Common
}

// TransGlobalPage 转换成全局名单分页需要的数据
func (t *TrustDirTrans) TransGlobalPage(total int64, models []*trust_dir_models.GlobalModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {

		for _, model := range models {
			data := make(map[string]interface{})
			data["id"] = fmt.Sprintf("%d", model.Id)
			data["name"] = model.Name
			data["dir"] = model.Dir
			data["platform_type"] = model.PlatformType
			data["platform_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("platform_type", model.PlatformType, enum_software.EnumMap)
			data["type"] = model.Type
			data["is_extend"] = model.IsExtend
			data["is_extend_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("is_extend", model.IsExtend, enum_tamper.EnumMapInt)
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
