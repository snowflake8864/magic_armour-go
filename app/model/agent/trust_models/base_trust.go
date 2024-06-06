package trust_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	OrgId        int64
	AgentId      string
	Id           int64
	PlatformType int
	DeptId       int64
	DeptIds      []int64
	OpenSwitch   *int
	Name         string
}

type baseTrust struct {
}

// conditions2Where 查询条件转化
func (m *baseTrust) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if condition.PlatformType > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "platform_type",
			Val: condition.PlatformType,
		})
	}
	if condition.Id > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}
	if condition.DeptId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dept_id",
			Val: condition.DeptId,
		})
	}

	if len(condition.DeptIds) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dept_id in",
			Val: condition.DeptIds,
		})
	}

	if condition.OpenSwitch != nil && *condition.OpenSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "open_switch",
			Val: *condition.OpenSwitch,
		})
	}
	if !string_helper.IsEmpty(condition.Name) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "name like",
			Val: "%" + condition.Name + "%",
		})
	}

	return wheres
}
