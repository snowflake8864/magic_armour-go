package user_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId          int64
	Id             int64
	CreateAtGte    string
	CreateAtLte    string
	LastLoginAtGte string
	LastLoginAtLte string
}

type baseUser struct {
}

// conditions2Where 查询条件转化
func (crd *baseUser) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.Id > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}

	if !string_helper.IsEmpty(condition.CreateAtGte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "created_at >=",
			Val: condition.CreateAtGte,
		})
	}
	if !string_helper.IsEmpty(condition.CreateAtLte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "created_at <=",
			Val: condition.CreateAtLte,
		})
	}
	if !string_helper.IsEmpty(condition.LastLoginAtGte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "last_login_at >=",
			Val: condition.LastLoginAtGte,
		})
	}
	if !string_helper.IsEmpty(condition.LastLoginAtLte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "last_login_at <=",
			Val: condition.LastLoginAtLte,
		})
	}

	return wheres
}

func (crd *baseUser) Columns() []string {
	return []string{
		"id",
		"username",
		"org_id",
		"role_id",
		"email",
		"phone",
		"realname",
		"position",
		"intro",
		"pw_error_num",
		"last_login_at",
		"is_root",
		"depts",
		"status",
		"create_id",
		"create_name",
		"created_at",
	}
}
