package plug_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId      int64
	AgentId    string
	Ip         string
	OpenSwitch *int
}

type basePlug struct {
}

// conditions2Where 查询条件转化
func (m *basePlug) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if !string_helper.IsEmpty(condition.Ip) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "ip",
			Val: condition.Ip,
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
