package org_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_user"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

type User struct {
	transformer.Common
}

func (t *User) pageGeneral(model *user_models.UserModel, departments map[int64]*department_models.DepartmentModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["id"] = model.Id
	data["username"] = model.Username
	data["email"] = model.Email
	data["phone"] = model.Phone
	data["realname"] = model.Realname
	data["position"] = model.Position
	data["intro"] = model.Intro
	data["last_login_at"] = model.LastLoginAt
	data["is_root"] = model.IsRoot
	depts := json_helper.ToSliceInt64(model.Depts)
	var deptNames []interface{}
	data["depts"] = depts
	if len(depts) > 0 {
		for _, did := range depts {
			if deptModel, ok := departments[did]; ok && !string_helper.IsEmpty(deptModel.FullName) {
				deptNames = append(deptNames, json_helper.ToSliceString(deptModel.FullName))
			}
		}
	}
	data["depts"] = depts
	data["dept_fullname"] = deptNames

	data["status"] = model.Status
	data["status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("status", model.Status, enum_user.EnumMap)
	data["is_root"] = model.IsRoot
	data["create_id"] = model.CreateId
	data["create_name"] = model.CreateName
	data["created_at"] = model.CreatedAt

	return data
}

// TransPage 转换成管理页需要的数据
func (t *User) TransPage(total int64, models []*user_models.UserModel, departments map[int64]*department_models.DepartmentModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model, departments)

			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransInfo 转换成管理页需要的数据
func (t *User) TransInfo(model *user_models.UserModel, departments map[int64]*department_models.DepartmentModel) map[string]interface{} {
	return t.pageGeneral(model, departments)
}
