package proc_models

import "goskeleton/app/utils/helpers/string_helper"

type SearchConditions struct {
	OrgId        int64
	AgentId      string
	AgentIds     []string
	ProcHash     string
	DllHash      string
	Id           string
	Status       int
	SearchStatus int
	Dir          string
	GteCreateAt  string
	LteCreateAt  string
}

type baseProc struct {
}

// conditions2Where 查询条件转化
func (m *baseProc) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if len(search.AgentIds) >= 0 {
		wheres["agent_id in"] = search.AgentIds
	}
	if !string_helper.IsEmpty(search.Id) {
		wheres["id"] = search.Id
	}

	if search.Status > 0 {
		wheres["status"] = search.Status
	}
	return wheres
}
