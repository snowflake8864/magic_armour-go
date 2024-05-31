package port_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
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
func (m *basePort) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if !string_helper.IsEmpty(condition.Id) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}
	if !string_helper.IsEmpty(condition.Port) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "port",
			Val: condition.Port,
		})
	}
	if condition.OpenSwitch != nil && *condition.OpenSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "open_switch",
			Val: *condition.OpenSwitch,
		})
	}

	return wheres
}
