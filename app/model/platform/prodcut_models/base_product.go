package prodcut_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	Id        int64
	GrantType int
	Title     string
	IsPublish int
}
type baseProduct struct {
}

// search2Where 查询条件转化
func (crd *baseProduct) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.Id > 0 {
		wheres["id"] = search.Id
	}
	if search.GrantType > 0 && search.GrantType != enums.AllDefaultInt {
		wheres["grant_type"] = search.GrantType
	}
	if !string_helper.IsEmpty(search.Title) {
		wheres["title like"] = "%" + search.Title + "%"
	}
	if search.IsPublish > 0 && search.IsPublish != enums.AllDefaultInt {
		wheres["is_publish"] = search.IsPublish
	}
	return wheres
}
