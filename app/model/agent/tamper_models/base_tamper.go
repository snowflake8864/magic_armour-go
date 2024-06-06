package tamper_models

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	OrgId        int64
	AgentId      string
	AgentIds     []string
	Id           int64
	PlatformType int
	DeptId       int64
	DeptIds      []int64
	OpenSwitch   *int
	DirFull      string
}

// Calc 统计
type Calc struct {
	AgentId     string `json:"agent_id"`
	TotalNumber int    `json:"total_number"`
	OnNumber    int    `json:"on_number"`
}

type baseTamper struct {
}

// ColumnCalc 字段：统计
func (m *baseTamper) ColumnCalc(switchOn bool) string {
	var columns []string
	columns = append(columns, "agent_id AS agent_id")
	columns = append(columns, "COUNT(id) AS total_number")
	if switchOn {
		columns = append(columns, fmt.Sprintf("SUM(IF(open_switch=%d, 1, 0)) AS on_number", enums.SwitchOn))
	}
	return strings.Join(columns, ", ")
}

// conditions2Where 查询条件转化
func (m *baseTamper) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if len(condition.AgentIds) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "agent_id in",
			Val: condition.AgentIds,
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

	if !string_helper.IsEmpty(condition.DirFull) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dir",
			Val: condition.DirFull,
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
