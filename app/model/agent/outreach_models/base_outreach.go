package outreach_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId      int64
	Ip         string
	Type       int
	OpenSwitch int
}

type baseOutreach struct {
}

// conditions2Where 查询条件转化
func (m *baseOutreach) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})

	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.Ip) {
		wheres["ip"] = search.Ip
	}
	if search.Type > 0 {
		wheres["type"] = search.Type
	}
	if search.OpenSwitch != enums.AllDefaultInt {
		wheres["open_switch"] = search.OpenSwitch
	}
	return wheres
}
