package syslog_models

import (
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	OrgId       int64
	AgentId     string
	CreateAtGte string
	CreateAtLte string
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
func (m *baseSyslog) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if !string_helper.IsEmpty(search.SourceIp) {
		wheres["source_ip"] = search.SourceIp
	}
	if search.LoginStatus != nil {
		wheres["handle_status"] = *search.LoginStatus
	}
	if !string_helper.IsEmpty(search.CreateAtGte) {
		wheres["created_at >="] = search.CreateAtGte
	}
	if !string_helper.IsEmpty(search.CreateAtLte) {
		wheres["created_at <="] = search.CreateAtLte
	}

	return wheres
}

// ColumnCountLogin 字段：累计登录日志次数
func (m *baseSyslog) ColumnCountLogin() string {
	columns := []string{"COUNT(id) AS total_number"}

	return strings.Join(columns, ", ")
}
