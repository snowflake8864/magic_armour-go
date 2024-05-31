package agent_monitor_models

import "goskeleton/app/utils/helpers/string_helper"

type SearchConditions struct {
	OrgId      int64
	AgentId    string
	MonitorDay string
	CreatedAt  string
}
type baseMonitor struct {
}

// search2Where 查询条件转化
func (crd *baseMonitor) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if !string_helper.IsEmpty(search.MonitorDay) {
		wheres["monitor_day"] = search.MonitorDay
	}
	if !string_helper.IsEmpty(search.CreatedAt) {
		wheres["created_at"] = search.CreatedAt
	}

	return wheres
}
