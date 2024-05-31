package agent_monitor_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId      int64
	AgentId    string
	MonitorDay string
	CreatedAt  string
}
type baseMonitor struct {
}

// search2Where 查询条件转化
func (crd *baseMonitor) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}
	if !string_helper.IsEmpty(condition.AgentId) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "agent_id",
			Val: condition.AgentId,
		})
	}
	if !string_helper.IsEmpty(condition.MonitorDay) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "monitor_day",
			Val: condition.MonitorDay,
		})
	}
	if !string_helper.IsEmpty(condition.CreatedAt) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "created_at",
			Val: condition.CreatedAt,
		})
	}

	return wheres
}
