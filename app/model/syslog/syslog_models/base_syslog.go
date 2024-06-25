package syslog_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	OrgId       int64
	AgentId     string
	AgentIds    []string
	Ids         []int64
	CreateAtGte string
	CreateAtLte string
	ProcHash    string
	DomainName  string
	SourceIp    string
	LoginStatus *int
}

// SyslogCalc 告警统计
type SyslogCalc struct {
	OrgId       int64  `json:"org_id"`
	AgentId     string `json:"agent_id"`
	TotalNumber int    `json:"total_number"`
}

type baseSyslog struct {
}

// conditions2Where 查询条件转化
func (m *baseSyslog) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if len(condition.Ids) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id in",
			Val: condition.Ids,
		})
	}
	if !string_helper.IsEmpty(condition.SourceIp) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "source_ip",
			Val: condition.SourceIp,
		})
	}
	if condition.LoginStatus != nil {
		wheres = append(wheres, model_comm.Conditions{
			Key: "handle_status",
			Val: *condition.LoginStatus,
		})
	}
	if !string_helper.IsEmpty(condition.CreateAtGte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "created_at >=",
			Val: condition.CreateAtGte,
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

// ColumnCountLogin 字段：累计登录日志次数
func (m *baseSyslog) ColumnCountLogin() string {
	columns := []string{"COUNT(id) AS total_number"}

	return strings.Join(columns, ", ")
}
