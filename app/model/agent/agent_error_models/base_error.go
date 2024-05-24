package agent_error_models

import (
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId   int64
	AgentId string
	Id      int64
}

type baseError struct {
}

// conditions2Where 查询条件转化
func (m *baseError) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if search.Id > 0 {
		wheres["id"] = search.Id
	}
	return wheres
}
