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
	ProcId      int
	ProcDir     string
	ProcParam   string
	ParentHash  string
	ParentId    int
	ParentDir   string
	ParentParam string
	DomainName  string
	SourceIp    string
	SourcePort  int
	TargetIp    string
	TargetPort  int
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
	if !string_helper.IsEmpty(condition.ProcHash) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "proc_hash",
			Val: condition.ProcHash,
		})
	}
	if condition.ProcId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "proc_id",
			Val: condition.ProcId,
		})
	}
	if !string_helper.IsEmpty(condition.ProcDir) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "proc_dir like",
			Val: "%" + condition.ProcDir + "%",
		})
	}
	if !string_helper.IsEmpty(condition.ProcParam) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "proc_param like",
			Val: "%" + condition.ProcParam + "%",
		})
	}
	if !string_helper.IsEmpty(condition.ParentHash) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "parent_hash",
			Val: condition.ParentHash,
		})
	}
	if condition.ParentId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "parent_id",
			Val: condition.ParentId,
		})
	}
	if !string_helper.IsEmpty(condition.ParentDir) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "parent_dir like",
			Val: "%" + condition.ParentDir + "%",
		})
	}
	if !string_helper.IsEmpty(condition.ParentParam) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "parent_param like",
			Val: "%" + condition.ParentParam + "%",
		})
	}
	if !string_helper.IsEmpty(condition.DomainName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "domain_name",
			Val: condition.DomainName,
		})
	}
	if !string_helper.IsEmpty(condition.SourceIp) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "source_ip",
			Val: condition.SourceIp,
		})
	}
	if condition.SourcePort > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "source_port",
			Val: condition.SourcePort,
		})
	}
	if !string_helper.IsEmpty(condition.TargetIp) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "target_ip",
			Val: condition.TargetIp,
		})
	}
	if condition.TargetPort > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "target_port",
			Val: condition.TargetPort,
		})
	}
	if condition.LoginStatus != nil {
		wheres = append(wheres, model_comm.Conditions{
			Key: "status",
			Val: *condition.LoginStatus,
		})
	}

	return wheres
}

// ColumnCountLogin 字段：累计登录日志次数
func (m *baseSyslog) ColumnCountLogin() string {
	columns := []string{"COUNT(id) AS total_number"}

	return strings.Join(columns, ", ")
}
