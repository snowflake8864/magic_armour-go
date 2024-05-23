package organization_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId      int64
	OrgType    int
	OraName    string
	LogOraName string
	UserId     int64
	UserName   string
	UserPhone  string
	StatusList []int
}

type baseOrg struct {
	OrgId       int64
	AgentId     string
	StatMonth   int
	StatDay     string
	GteCreateAt string
	LteCreateAt string
}

// search2Where 查询条件转化
func (crd *baseOrg) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if search.OrgType > 0 && search.OrgType != enums.AllDefaultInt {
		wheres["type"] = search.OrgType
	}
	if !string_helper.IsEmpty(search.OraName) {
		wheres["name like"] = "%" + search.OraName + "%"
	}
	if !string_helper.IsEmpty(search.LogOraName) {
		wheres["org_name like"] = "%" + search.LogOraName + "%"
	}
	if search.UserId > 0 {
		wheres["user_id"] = search.UserId
	}
	if !string_helper.IsEmpty(search.UserName) {
		wheres["user_name"] = search.UserName
	}
	if !string_helper.IsEmpty(search.UserPhone) {
		wheres["user_phone"] = search.UserPhone
	}
	if search.StatusList != nil && len(search.StatusList) > 0 {
		wheres["status in"] = search.StatusList
	}
	return wheres
}
