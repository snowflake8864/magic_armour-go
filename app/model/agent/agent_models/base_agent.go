package agent_models

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

// AgentCalc 终端统计
type AgentCalc struct {
	OrgId               int64 `json:"org_id"`
	TotalNumber         int   `json:"total_number"`
	OnlineNumber        int   `json:"online_number"`
	OfflineNumber       int   `json:"offline_number"`
	SafeNumber          int   `json:"safe_number"`
	UnsafeNumber        int   `json:"unsafe_number"`
	WindowsNumber       int   `json:"windows_number"`
	WindowsOnlineNumber int   `json:"windows_online_number"`
	LinuxNumber         int   `json:"linux_number"`
	LinuxOnlineNumber   int   `json:"linux_online_number"`
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
	GteCreateAt  string
	LteCreateAt  string
}

type baseAgent struct {
}

// conditions2Where 查询条件转化
func (m *baseAgent) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if search.OrgId > 0 {
		wheres["org_id"] = search.OrgId
	}
	if !string_helper.IsEmpty(search.AgentId) {
		wheres["id"] = search.AgentId
	}
	if search.AgentIds != nil && len(search.AgentIds) > 0 {
		wheres["id in"] = search.AgentIds
	}
	if !string_helper.IsEmpty(search.HostName) {
		wheres["host_name"] = search.HostName
	}
	if !string_helper.IsEmpty(search.NetIp) {
		wheres["net_ip"] = search.NetIp
	}
	if !string_helper.IsEmpty(search.AgentIp) {
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
	if !string_helper.IsEmpty(search.ContactName) {
		wheres["contact_name"] = search.ContactName
	}
	if !string_helper.IsEmpty(search.Alias) {
		wheres["alias"] = search.Alias
	}
	if search.DeptIds != nil && len(search.DeptIds) > 0 {
		wheres["dept_id in"] = search.DeptIds
	}
	if !string_helper.IsEmpty(search.GteCreateAt) {
		wheres["created_at >="] = search.GteCreateAt
	}
	if !string_helper.IsEmpty(search.LteCreateAt) {
		wheres["created_at <="] = search.LteCreateAt
	}
	return wheres
}

// ColumnSimplify 字段: 最简单的
func (m *baseAgent) ColumnSimplify() string {
	return "org_id, id"
}

// ColumnBasic 字段: 基础基础
func (m *baseAgent) ColumnBasic(needId bool) string {
	column := ""
	if needId {
		column = "org_id, id,"
	}
	column = column +
		"platform_type, os_type, host_name, mac_id, agent_ip, net_ip, os_name, mem_desc, " +
		"hd_desc, cpu_desc, network_card, agent_start_at, os_start_at, ver," +
		"is_online, online_at, live_at, is_safe, dept_id,user_id," +
		"agent_type,alias,contact_name,contact_tel, position, remark"
	return column
}

// ColumnConfig 字段:全部配置相关
func (m *baseAgent) ColumnConfig(needId, grant bool) string {
	column := ""
	if needId {
		column = "org_id, id,"
	}
	if grant {
		column += "grant_config, grant_order_no,grant_expire_at,grant_status,"
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
		column = "org_id, id, os_name,"
	}
	column = column +
		"platform_type, os_type, process_switch, process_mode, dll_switch, dll_mode, auto_white_switch, alarm_ai_switch"
	return column
}

// ColumnSwitchAndModeForAlarm 字段:告警需求的开关和模式
func (m *baseAgent) ColumnSwitchAndModeForAlarm(needId bool) string {
	column := ""
	if needId {
		column = "org_id, id,"
	}
	column = column +
		"platform_type, os_type, process_switch, process_mode, dll_switch, dll_mode, tamper_switch, tamper_mode, " +
		"blackmail_switch, blackmail_mode, port_switch,plug_switch, usb_switch, usb_mode"
	return column
}

// ColumnCalcAgents 字段： 统计终端
func (m *baseAgent) ColumnCalcAgents(org, online, safe, win, winOnline bool) string {
	columns := []string{"COUNT(id) AS total_number"}
	if org {
		columns = append(columns, "org_id")
	}
	if online {
		columns = append(columns, "SUM(IF(is_online ="+int_helper.IntToStr(enums.IsYes)+", 1, 0)) AS online_number")
	}
	if safe {
		columns = append(columns, "SUM(IF(is_safe ="+int_helper.IntToStr(enums.IsYes)+", 1, 0)) AS safe_number")
	}
	if win {
		columns = append(columns, "SUM(IF(platform_type ="+int_helper.IntToStr(enum_software.PlatformTypeWindows)+", 1, 0)) AS windows_number")
	}
	if winOnline {
		columns = append(columns, fmt.Sprintf("SUM(IF(platform_type=%d AND is_online=%d, 1, 0)) AS windows_online_number", enum_software.PlatformTypeWindows, enums.IsYes))
	}
	return strings.Join(columns, ", ")
}

func (m *baseAgent) FillCalcData(data *AgentCalc) *AgentCalc {
	data.UnsafeNumber = data.TotalNumber - data.SafeNumber
	data.LinuxNumber = data.TotalNumber - data.WindowsNumber
	data.LinuxOnlineNumber = data.OnlineNumber - data.WindowsOnlineNumber

	return data
}
