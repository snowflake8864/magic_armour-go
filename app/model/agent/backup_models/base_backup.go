package backup_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	OrgId            int64
	AgentId          string
	AgentIds         []string
	Id               string
	SourcePath       string
	BackupStatus     int
	BackupStatusList []int
	RollbackStatus   int
}

// Calc 统计
type Calc struct {
	AgentId     string `json:"agent_id"`
	TotalNumber int    `json:"total_number"`
}

type baseBackup struct {
}

// ColumnCalc 字段：统计
func (m *baseBackup) ColumnCalc() string {
	var columns []string
	columns = append(columns, "agent_id AS agent_id")
	columns = append(columns, "COUNT(id) AS total_number")

	return strings.Join(columns, ", ")
}

// conditions2Where 查询条件转化
func (m *baseBackup) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if !string_helper.IsEmpty(condition.Id) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.Id,
		})
	}

	if !string_helper.IsEmpty(condition.SourcePath) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "source_path",
			Val: condition.SourcePath,
		})
	}

	if condition.BackupStatus > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "backup_status",
			Val: condition.BackupStatus,
		})
	}
	if len(condition.BackupStatusList) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "backup_status in",
			Val: condition.BackupStatusList,
		})
	}

	if condition.RollbackStatus > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "rollback_status",
			Val: condition.RollbackStatus,
		})
	}

	return wheres
}
