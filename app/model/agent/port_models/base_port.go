package port_models

import (
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId      int64
	AgentId    string
	Id         string
	Port       string
	OpenSwitch *int
}

type basePort struct {
}

// conditions2Where 查询条件转化
func (m *basePort) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if !string_helper.IsEmpty(search.Id) {
		wheres["id"] = search.Id
	}
	if !string_helper.IsEmpty(search.Port) {
		wheres["port"] = search.Port
	}
	if search.OpenSwitch != nil {
		wheres["open_switch"] = &search.OpenSwitch
	}

	return wheres
}
