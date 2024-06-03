package proc_models

import (
	"fmt"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

// StatusCalc 状态统计
type StatusCalc struct {
	AgentId       string `json:"agent_id"`
	TotalNumber   int    `json:"total_number"`
	UnsetNumber   int    `json:"unset_number"`
	UnknownNumber int    `json:"unknown_number"`
	WhiteNumber   int    `json:"white_number"`
	BlackNumber   int    `json:"black_number"`
}

// RepeatCalc 重复hash统计
type RepeatCalc struct {
	Hash        int64 `json:"hash"`
	TotalNumber int   `json:"total_number"`
}

type SearchConditions struct {
	OrgId           int64
	AgentId         string
	AgentIds        []string
	ProcHash        string
	DllHash         string
	Id              string
	Status          int
	StatusList      []int
	Dir             string
	Dirs            []string
	SyncStatusAtGte int64
}

type baseProc struct {
}

// conditions2Where 查询条件转化
func (m *baseProc) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if !string_helper.IsEmpty(condition.ProcHash) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "proc_hash",
			Val: condition.ProcHash,
		})
	}
	if !string_helper.IsEmpty(condition.DllHash) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dll_hash",
			Val: condition.DllHash,
		})
	}
	if !string_helper.IsEmpty(condition.Id) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}
	if !string_helper.IsEmpty(condition.Dir) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dir like",
			Val: "%" + condition.Dir + "%",
		})

	}
	if len(condition.Dirs) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dir",
			Val: condition.Dirs,
		})
	}
	if condition.Status > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "status",
			Val: condition.Status,
		})
	}
	if len(condition.StatusList) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "status in",
			Val: condition.StatusList,
		})
	}
	if condition.SyncStatusAtGte > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "sync_status_at >=",
			Val: condition.SyncStatusAtGte,
		})
	}
	return wheres
}

// ColumnStatusCalc 字段： 统计状态
func (m *baseProc) ColumnStatusCalc(total, unset, unknown, white, black bool) string {
	var columns []string
	columns = append(columns, "agent_id")
	if !total && !unset && !unknown && !white && !black {
		columns = append(columns, "COUNT(id) AS total_number")
		columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS unset_number", enum_process.ProcStatusUnset))
		columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS unknown_number", enum_process.ProcStatusUnknown))
		columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS white_number", enum_process.ProcStatusWhite))
		columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS black_number", enum_process.ProcStatusBlack))
	}

	if total {
		columns = append(columns, "COUNT(id) AS total_number")
	}
	if unset {
		columns = append(columns, "COUNT(id) AS unset_number")
	}
	if unknown {
		columns = append(columns, "COUNT(id) AS unknown_number")
	}
	if white {
		columns = append(columns, "COUNT(id) AS white_number")
	}
	if black {
		columns = append(columns, "COUNT(id) AS black_number")
	}

	return strings.Join(columns, ", ")
}

// ColumnRepeatCalc 字段：重复hash统计
func (m *baseProc) ColumnRepeatCalc(process bool) string {
	var columns []string
	if process {
		columns = append(columns, "proc_hash AS hash")
	} else {
		columns = append(columns, "dll_hash AS hash")
	}

	columns = append(columns, "COUNT(id) AS total_number")
	return strings.Join(columns, ", ")
}
