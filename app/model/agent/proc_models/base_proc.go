package proc_models

import (
	"fmt"
	"goskeleton/app/global/enums/enum_process"
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
	Dir             string
	Dirs            []string
	GteSyncStatusAt string
}

type baseProc struct {
}

// conditions2Where 查询条件转化
func (m *baseProc) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if len(search.AgentIds) > 0 {
		wheres["agent_id in"] = search.AgentIds
	}
	if !string_helper.IsEmpty(search.ProcHash) {
		wheres["proc_hash"] = search.ProcHash
	}
	if !string_helper.IsEmpty(search.DllHash) {
		wheres["dll_hash"] = search.DllHash
	}
	if !string_helper.IsEmpty(search.Id) {
		wheres["id"] = search.Id
	}
	if !string_helper.IsEmpty(search.Dir) {
		wheres["dir like"] = "%" + search.Dir + "%"
	}
	if len(search.Dirs) > 0 {
		wheres["dir in"] = search.Dirs
	}
	if search.Status > 0 {
		wheres["status"] = search.Status
	}
	if !string_helper.IsEmpty(search.GteSyncStatusAt) {
		wheres["sync_status_at >="] = search.GteSyncStatusAt
	}
	return wheres
}

// ColumnStatusCalc 字段： 统计状态
func (m *baseProc) ColumnStatusCalc() string {
	var columns []string
	columns = append(columns, "agent_id")
	columns = append(columns, "COUNT(id) AS total_number")
	columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS unset_number", enum_process.ProcStatusUnset))
	columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS unknown_number", enum_process.ProcStatusUnknown))
	columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS white_number", enum_process.ProcStatusWhite))
	columns = append(columns, fmt.Sprintf("SUM(IF(status=%d, 1, 0)) AS black_number", enum_process.ProcStatusBlack))
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
