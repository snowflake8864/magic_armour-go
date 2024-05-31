package user_models

import "goskeleton/app/model/model_comm"

type SearchConditions struct {
	OrgId int64
	Id    int64
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
