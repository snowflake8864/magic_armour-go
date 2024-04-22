package organize_service

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers"
)

type SearchConditions struct {
	OrgId      int64
	OrgType    int
	OraName    string
	UserId     int64
	UserName   string
	UserPhone  string
	StatusList []int
}

type OrgBase struct {
}

// search2Where 查询条件转化
func (crd *OrgBase) search2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if search.OrgType > 0 && search.OrgType != enums.AllDefaultInt {
		wheres["type"] = search.OrgType
	}
	if !helpers.StringIsEmpty(search.OraName) {
		wheres["name like"] = "%" + search.OraName + "%"
	}
	if search.UserId > 0 {
		wheres["user_id"] = search.UserId
	}
	if !helpers.StringIsEmpty(search.UserName) {
		wheres["user_name"] = search.UserName
	}
	if !helpers.StringIsEmpty(search.UserPhone) {
		wheres["user_phone"] = search.UserPhone
	}
	if search.StatusList != nil && len(search.StatusList) > 0 {
		wheres["status in"] = search.StatusList
	}
	return wheres
}
