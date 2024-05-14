package agent_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/utils/helpers"
	"strings"
)

// AgentCalc 终端统计
type AgentCalc struct {
	OrgId               int64 `json:"org_id"`
	TotalNumber         int   `json:"total_number"`
	OnlineNumber        int   `json:"online_number"`
	SafeNumber          int   `json:"safe_number"`
	WindowsNumber       int   `json:"windows_number"`
	WindowsOnlineNumber int   `json:"windows_online_number"`
}

type SearchConditions struct {
	OrgId        int64
	AgentId      string
	AgentIds     []string
	HostName     string
	AgentIp      string
	NetIp        string
	PlatformType int
	IsOnline     int
	IsSafe       int
	Alias        string
	ContactName  string
	DeptIds      []int64
	CreateAtGte  string
	CreateAtLte  string
}

type baseAgent struct {
}

// conditions2Where 查询条件转化
func (m *baseAgent) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !helpers.StringIsEmpty(search.AgentId) {
		wheres["agent_id"] = search.AgentId
	}
	if search.AgentIds != nil && len(search.AgentIds) > 0 {
		wheres["agent_id in"] = search.AgentIds
	}
	if !helpers.StringIsEmpty(search.HostName) {
		wheres["host_name"] = search.HostName
	}
	if !helpers.StringIsEmpty(search.NetIp) {
		wheres["net_ip"] = search.NetIp
	}
	if !helpers.StringIsEmpty(search.AgentIp) {
		wheres["agent_ip like"] = "%" + search.AgentIp + "%"
	}
	if search.PlatformType > 0 {
		wheres["platform_type"] = search.PlatformType
	}
	if search.IsOnline != enums.AllDefaultInt {
		wheres["is_online"] = search.IsOnline
	}
	if search.IsSafe != enums.AllDefaultInt {
		wheres["is_safe"] = search.IsSafe
	}
	if !helpers.StringIsEmpty(search.ContactName) {
		wheres["contact_name"] = search.ContactName
	}
	if !helpers.StringIsEmpty(search.Alias) {
		wheres["alias"] = search.Alias
	}
	if search.DeptIds != nil && len(search.DeptIds) > 0 {
		wheres["dept_id in"] = search.DeptIds
	}
	if !helpers.StringIsEmpty(search.CreateAtGte) {
		wheres["created_at >="] = search.CreateAtGte
	}
	if !helpers.StringIsEmpty(search.CreateAtLte) {
		wheres["created_at <="] = search.CreateAtLte
	}
	return wheres
}

// ColumnConfig 字段:全部配置相关
func (m *baseAgent) ColumnConfig(needId bool) string {
	column := ""
	if needId {
		column = "org_id, agent_id"
	}
	column = column +
		"platform_type, os_type, process_switch, process_mode, dll_switch, dll_mode, tamper_switch, tamper_mode, " +
		"blackmail_switch, blackmail_mode, port_switch, plug_switch, usb_switch, usb_mode," +
		"syslog_dns_switch, syslog_outer_switch, syslog_inner_switch, syslog_process_switch, syslog_login_switch," +
		"outreach_alarm_switch, self_protect_switch, kernel_upgrade_switch, auto_white_switch, alarm_ai_switch," +
		"login_solve_switch, hardware_switch, hardware_time, common_time,fast_time"
	return column
}

// ColumnProc 字段:全部可信防护相关
func (m *baseAgent) ColumnProc(needId bool) string {
	column := ""
	if needId {
		column = "org_id, agent_id, os_name"
	}
	column = column +
		"platform_type, os_type, process_switch, process_mode, dll_switch, dll_mode, auto_white_switch, alarm_ai_switch"
	return column
}

// ColumnSwitchAndModeForAlarm 字段:告警需求的开关和模式
func (m *baseAgent) ColumnSwitchAndModeForAlarm(needId bool) string {
	column := ""
	if needId {
		column = "org_id, agent_id,"
	}
	column = column +
		"platform_type, os_type, process_switch, process_mode, dll_switch, dll_mode, tamper_switch, tamper_mode, " +
		"blackmail_switch, blackmail_mode, port_switch,plug_switch, usb_switch, usb_mode"
	return column
}

// ColumnCalcAgents 字段： 统计终端
func (m *baseAgent) ColumnCalcAgents(org, online, safe, win, winOnline bool) string {
	columns := []string{"COUNT(agent_id) AS total_number"}
	if org {
		columns = append(columns, "org_id")
	}
	if online {
		columns = append(columns, "SUM(IF(is_online ="+helpers.IntToStr(enums.IsYes)+", 1, 0)) AS online_number")
	}
	if safe {
		columns = append(columns, "SUM(IF(is_safe ="+helpers.IntToStr(enums.IsYes)+", 1, 0)) AS safe_number")
	}
	if win {
		columns = append(columns, "SUM(IF(platform_type ="+helpers.IntToStr(enum_software.PlatformTypeWindows)+", 1, 0) AS windows_number")
	}
	if winOnline {
		columns = append(columns, "SUM(IF(platform_type ="+helpers.IntToStr(enum_software.PlatformTypeWindows)+", 1, 0)"+
			" AND is_online ="+helpers.IntToStr(enums.IsYes)+", 1, 0)) AS  windows_online_number")
	}

	return strings.Join(columns, ", ")
}
