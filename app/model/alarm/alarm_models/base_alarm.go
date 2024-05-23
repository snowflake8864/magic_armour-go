package alarm_models

import (
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	OrgId        int64
	AgentId      string
	AttackIp     string
	HandleStatus int
	GteCreateAt  string
	LteCreateAt  string
}

// AlarmCalc 告警统计
type AlarmCalc struct {
	OrgId           int64  `json:"org_id"`
	AgentId         string `json:"agent_id"`
	TotalNumber     int    `json:"total_number"`
	HandledNumber   int    `json:"handled_number"`
	UnhandledNumber int    `json:"unhandled_number"`
}

type baseAlarm struct {
}

// conditions2Where 查询条件转化
func (m *baseAlarm) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if !string_helper.IsEmpty(search.AttackIp) {
		wheres["attack_ip"] = search.AttackIp
	}
	if search.HandleStatus > 0 {
		wheres["handle_status"] = search.HandleStatus
	}
	if !string_helper.IsEmpty(search.GteCreateAt) {
		wheres["created_at >="] = search.GteCreateAt
	}
	if !string_helper.IsEmpty(search.LteCreateAt) {
		wheres["created_at <="] = search.LteCreateAt
	}

	return wheres
}

// ColumnSumWeight 字段：累计端口虚开告警权重
func (m *baseAlarm) ColumnSumWeight() string {
	columns := []string{"SUM(weight*attack_total) AS total_number"}

	return strings.Join(columns, ", ")
}

// ColumnCalc 字段：统计字段
func (m *baseAlarm) ColumnCalc(org, agent, handled, unhandled bool) string {
	columns := []string{"COUNT(id) AS total_number"}
	if org {
		columns = append(columns, "org_id")
	}
	if agent {
		columns = append(columns, "agent_id")
	}
	if handled {
		columns = append(columns, "SUM(IF(handle_status ="+int_helper.IntToStr(enum_alarm.HandleStatusYes)+", 1, 0)) AS handled_number")
	}
	if unhandled {
		columns = append(columns, "SUM(IF(handle_status ="+int_helper.IntToStr(enum_alarm.HandleStatusNot)+", 1, 0)) AS unhandled_number")
	}

	return strings.Join(columns, ", ")
}
