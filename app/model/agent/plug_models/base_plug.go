package plug_models

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	OrgId      int64
	AgentId    string
	AgentIds   []string
	Id         int64
	DeptId     int64
	DeptIds    []int64
	Ip         string
	IpType     int
	OpenSwitch *int
	CreateId   int64
}

// Calc 统计
type Calc struct {
	AgentId     string `json:"agent_id"`
	TotalNumber int    `json:"total_number"`
	OnNumber    int    `json:"on_number"`
}

type basePlug struct {
}

// ColumnCalc 字段：统计
func (m *basePlug) ColumnCalc(switchOn bool) string {
	var columns []string
	columns = append(columns, "agent_id AS agent_id")
	columns = append(columns, "COUNT(id) AS total_number")
	if switchOn {
		columns = append(columns, fmt.Sprintf("SUM(IF(open_switch=%d, 1, 0)) AS on_number", enums.SwitchOn))
	}
	return strings.Join(columns, ", ")
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
	if len(condition.AgentIds) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "agent_id in",
			Val: condition.AgentIds,
		})
	}
	if condition.Id > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}
	if !string_helper.IsEmpty(condition.Ip) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "ip",
			Val: condition.Ip,
		})
	}
	if condition.IpType > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "type",
			Val: condition.IpType,
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
	if condition.CreateId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "create_id",
			Val: condition.CreateId,
		})
	}
	return wheres
}
