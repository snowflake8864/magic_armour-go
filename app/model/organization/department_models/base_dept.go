package department_models

import "goskeleton/app/utils/helpers/string_helper"

type SearchConditions struct {
	OrgId    int64
	Id       int64
	Ids      []int64
	FullCode string
}

type baseDept struct {
}

// conditions2Where 查询条件转化
func (crd *baseDept) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if search.Id > 0 {
		wheres["id"] = search.Id
	}
	if len(search.Ids) > 0 {
		wheres["id in"] = search.Ids
	}
	if !string_helper.IsEmpty(search.FullCode) {
		wheres["id findInSetLater"] = search.FullCode
	}
	return wheres
}
