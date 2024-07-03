package task_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId       int64
	AgentId     string
	Type        *int
	Types       []int
	Status      int
	CreateAtLte string
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
	if condition.Type != nil && *condition.Type != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "type",
			Val: *condition.Type,
		})
	}
	if len(condition.Types) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "type in",
			Val: condition.Types,
		})
	}

	if condition.Status > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "status",
			Val: condition.Status,
		})
	}

	if !string_helper.IsEmpty(condition.CreateAtLte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "created_at <=",
			Val: condition.CreateAtLte,
		})
	}
	return wheres
}
