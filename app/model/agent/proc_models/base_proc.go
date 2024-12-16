package proc_models

import (
	"fmt"
	"goskeleton/app/global/enums"
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
	Hash        string `json:"hash"`
	TotalNumber int    `json:"total_number"`
}

// ProcessAndRunning 进程和运行中标
type ProcessAndRunning struct {
	Id                  string `json:"id"`
	AgentId             string `json:"agent_id"`
	ProcHash            string `json:"proc_hash"`
	Status              int    `json:"status"`
	SyncStatusAt        int64  `json:"sync_status_at"`
	GlobalStatus        int    `json:"global_status"`
	SelfStatus          int    `json:"self_status"`
	ThreatenStatus      int    `json:"threaten_status"`
	ThreatenLevel       int    `json:"threaten_level"`
	AutoWhiteSwitch     int    `json:"auto_white_switch"`
	GrantThreatenSwitch int    `json:"grant_threaten_switch"`
	Introduce           string `json:"introduce"`
	Copyright           string `json:"copyright"`
	Signature           string `json:"signature"`
	CreatedAt           string `json:"created_at"`
	ProcId              int    `json:"proc_id"`
	ProcDir             string `json:"proc_dir"`
	ProcUser            string `json:"proc_user"`
	DllNumber           int    `json:"dll_number"`
}

// DllAndRunning 进程和运行中标
type DllAndRunning struct {
	Id                  string `json:"id"`
	AgentId             string `json:"agent_id"`
	DllHash             string `json:"dll_hash"`
	DllDir              string `json:"dll_dir"`
	Status              int    `json:"status"`
	SyncStatusAt        int64  `json:"sync_status_at"`
	GlobalStatus        int    `json:"global_status"`
	SelfStatus          int    `json:"self_status"`
	ThreatenStatus      int    `json:"threaten_status"`
	ThreatenLevel       int    `json:"threaten_level"`
	AutoWhiteSwitch     int    `json:"auto_white_switch"`
	GrantThreatenSwitch int    `json:"grant_threaten_switch"`
	Introduce           string `json:"introduce"`
	Copyright           string `json:"copyright"`
	Signature           string `json:"signature"`
	Attribute           string `json:"attribute"`
	CreatedAt           string `json:"created_at"`
}

type SearchConditions struct {
	OrgId               int64
	AgentId             string
	AgentIds            []string
	ProcHash            string
	ProcHashes          []string
	DllHash             string
	DllHashes           []string
	Id                  string
	Status              int
	StatusList          []int
	Dir                 string
	Dirs                []string
	ProcDir             string
	DllDir              string
	GlobalType          int
	DataFrom            int
	OpenSwitch          *int
	GrantThreatenSwitch *int
	AutoWhiteSwitch     *int
	ThreatenStatus      *int
	SyncStatusAtGte     int64
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
	if len(condition.ProcHashes) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "proc_hash in",
			Val: condition.ProcHashes,
		})
	}
	if !string_helper.IsEmpty(condition.DllHash) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dll_hash",
			Val: condition.DllHash,
		})
	}
	if len(condition.DllHashes) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dll_hash in",
			Val: condition.DllHashes,
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
	if condition.GlobalType > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "type",
			Val: condition.GlobalType,
		})
	}
	if condition.DataFrom > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "data_from",
			Val: condition.DataFrom,
		})
	}
	if condition.OpenSwitch != nil && *condition.OpenSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "open_switch",
			Val: *condition.OpenSwitch,
		})
	}
	if condition.AutoWhiteSwitch != nil && *condition.AutoWhiteSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "auto_white_switch",
			Val: *condition.AutoWhiteSwitch,
		})
	}
	if condition.GrantThreatenSwitch != nil && *condition.GrantThreatenSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "grant_threaten_switch",
			Val: *condition.GrantThreatenSwitch,
		})
	}
	if condition.ThreatenStatus != nil && *condition.ThreatenStatus != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "threaten_status",
			Val: *condition.ThreatenStatus,
		})
	}
	return wheres
}

// conditions2Where 查询条件转化
func (m *baseProc) conditions2WhereForRunning(condition SearchConditions, t1, t2 string) (wheres []model_comm.Conditions) {
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: t1 + "org_id",
			Val: condition.OrgId,
		})
	}
	if !string_helper.IsEmpty(condition.AgentId) {
		wheres = append(wheres, model_comm.Conditions{
			Key: t1 + "agent_id",
			Val: condition.AgentId,
		})
	}
	if !string_helper.IsEmpty(condition.ProcHash) {
		wheres = append(wheres, model_comm.Conditions{
			Key: t1 + "proc_hash",
			Val: condition.ProcHash,
		})
	}
	if !string_helper.IsEmpty(condition.DllHash) {
		wheres = append(wheres, model_comm.Conditions{
			Key: t1 + "dll_hash",
			Val: condition.DllHash,
		})
	}
	if !string_helper.IsEmpty(condition.ProcDir) {
		wheres = append(wheres, model_comm.Conditions{
			Key: t1 + "proc_dir like",
			Val: "%" + condition.ProcDir + "%",
		})
	}
	if !string_helper.IsEmpty(condition.DllDir) {
		wheres = append(wheres, model_comm.Conditions{
			Key: t1 + "dll_dir like",
			Val: "%" + condition.DllDir + "%",
		})
	}

	if condition.Status > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: t2 + "status",
			Val: condition.Status,
		})
	}
	return wheres
}

func (m *baseProc) getJoinStrForProcess(t1Alias, t2Name, t2Alias string) string {
	return "LEFT JOIN " + t2Name + " AS " + strings.Trim(t2Alias, ".") +
		" ON " + t1Alias + "org_id=" + t2Alias + "org_id" +
		" AND " + t1Alias + "agent_id=" + t2Alias + "agent_id " +
		" AND " + t1Alias + "proc_hash=" + t2Alias + "proc_hash"
}
func (m *baseProc) getJoinStrForDll(t1Alias, t2Name, t2Alias string) string {
	return "LEFT JOIN " + t2Name + " AS " + strings.Trim(t2Alias, ".") +
		" ON " + t1Alias + "org_id=" + t2Alias + "org_id" +
		" AND " + t1Alias + "agent_id=" + t2Alias + "agent_id " +
		" AND " + t1Alias + "dll_hash=" + t2Alias + "dll_hash"
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
