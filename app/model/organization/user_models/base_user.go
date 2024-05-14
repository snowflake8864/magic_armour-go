package user_models

type SearchConditions struct {
	OrgId int64
}

type baseUser struct {
	OrgId int64
}

// conditions2Where 查询条件转化
func (crd *baseUser) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	return wheres
}

func (m *UserModel) Columns() []string {
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
		"status",
		"create_id",
		"create_name",
		"created_at",
	}
}
