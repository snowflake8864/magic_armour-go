package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_user"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/model/organization/role_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/helpers/struct_helper"
)

type UserTrans struct {
	transformer.Common
	Lang string
}

func (t *UserTrans) pageGeneral(model *user_models.UserModel, departments map[int64]*department_models.DepartmentModel, roleMaps map[int64]*role_models.RoleModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)

	depts := json_helper.ToSliceInt64(model.Depts)
	data["depts"] = depts
	var deptNames []interface{}
	if len(depts) > 0 {
		for _, did := range depts {
			if deptModel, ok := departments[did]; ok && !string_helper.IsEmpty(deptModel.FullName) {
				deptNames = append(deptNames, json_helper.ToSliceString(deptModel.FullName))
			}
		}
	}
	data["dept_fullname"] = deptNames

	if t.Lang == consts.LanguageZh {
		data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_user.EnumMap)
	} else {
		data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status_en", model.Status, enum_user.EnumMap)

	}

	roleName := ""
	if roleModel, ok := roleMaps[model.RoleId]; ok && roleModel != nil {
		roleName = roleModel.Name
	}
	data["role_name"] = roleName

	return data
}

// TransPage 转换成管理页需要的数据
func (t *UserTrans) TransPage(total int64, models []*user_models.UserModel, departments map[int64]*department_models.DepartmentModel, roleMaps map[int64]*role_models.RoleModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model, departments, roleMaps)

			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransInfo 转换成管理页需要的数据
func (t *UserTrans) TransInfo(model *user_models.UserModel, departments map[int64]*department_models.DepartmentModel, roleMaps map[int64]*role_models.RoleModel) map[string]interface{} {
	return t.pageGeneral(model, departments, roleMaps)
}
