package alarm_models

import (
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type SearchConditions struct {
	OrgId        int64
	AgentId      string
	AgentIds     []string
	AttackIp     string
	HandleStatus int
	CreateAtGte  string
	CreateAtLte  string
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
func (m *baseAlarm) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
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
	if !string_helper.IsEmpty(condition.AttackIp) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "attack_ip",
			Val: condition.AttackIp,
		})
	}
	if condition.HandleStatus > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "handle_status",
			Val: condition.HandleStatus,
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

// ColumnSumWeight 字段：累计端口虚开告警权重
func (m *baseAlarm) ColumnSumWeight() string {
	columns := []string{"SUM(weight*attack_total) AS total_number"}

	return strings.Join(columns, ", ")
}

// ColumnCalc 字段：统计字段
func (m *baseAlarm) ColumnCalc(org, agent, unhandled, handled bool) string {
	columns := []string{"COUNT(id) AS total_number"}
	if org {
		columns = append(columns, "org_id")
	}
	if agent {
		columns = append(columns, "agent_id")
	}
	if unhandled {
		columns = append(columns, "SUM(IF(handle_status ="+int_helper.IntToStr(enum_alarm.HandleStatusNot)+", 1, 0)) AS unhandled_number")
	}
	if handled {
		columns = append(columns, "SUM(IF(handle_status ="+int_helper.IntToStr(enum_alarm.HandleStatusYes)+", 1, 0)) AS handled_number")
	}
	return strings.Join(columns, ", ")
}
