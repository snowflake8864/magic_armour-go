package task_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId   int64
	AgentId string
	Type    int
	Types   []int
}

type baseTask struct {
}

// conditions2Where 查询条件转化
func (m *baseTask) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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

	if len(condition.Types) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "type in",
			Val: condition.Types,
		})
	}
	if condition.Type > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "type",
			Val: condition.Type,
		})
	}
	return wheres
}
