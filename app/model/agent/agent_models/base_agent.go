package agent_models

import (
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/model/model_comm"
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
	OrgId               int64
	OrgIds              []int64
	AgentId             string
	AgentIds            []string
	HostName            string
	AgentIp             string
	NetIp               string
	PlatformType        int
	IsOnline            *int
	IsSafe              *int
	GrantStatus         int
	GrantExpireAtLte    string
	Alias               string
	ContactName         string
	DeptIds             []int64
	LiveAtLte           string
	CreateAtGte         string
	CreateAtLte         string
	ProcessSwitch       *int
	DllSwitch           *int
	TamperSwitch        *int
	BlackmailSwitch     *int
	PortSwitch          *int
	PlugSwitch          *int
	UsbSwitch           *int
	OutreachAlarmSwitch *int
}

type baseAgent struct {
}

// conditions2Where 查询条件转化
func (m *baseAgent) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	if condition.OrgId > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id",
			Val: condition.OrgId,
		})
	}
	if len(condition.OrgIds) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "org_id in",
			Val: condition.OrgIds,
		})
	}
	if !string_helper.IsEmpty(condition.AgentId) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id",
			Val: condition.AgentId,
		})
	}
	if len(condition.AgentIds) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "id in",
			Val: condition.AgentIds,
		})
	}
	if !string_helper.IsEmpty(condition.HostName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "host_name",
			Val: condition.HostName,
		})
	}
	if !string_helper.IsEmpty(condition.NetIp) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "net_ip",
			Val: condition.NetIp,
		})
	}
	if !string_helper.IsEmpty(condition.AgentIp) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "agent_ip like",
			Val: "%" + condition.AgentIp + "%",
		})
	}
	if condition.PlatformType > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "platform_type",
			Val: condition.PlatformType,
		})
	}
	if condition.IsOnline != nil && *condition.IsOnline != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_online",
			Val: *condition.IsOnline,
		})
	}
	if condition.IsSafe != nil && *condition.IsSafe != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "is_safe",
			Val: *condition.IsSafe,
		})
	}
	if condition.GrantStatus > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "grant_status",
			Val: condition.GrantStatus,
		})
	}
	if !string_helper.IsEmpty(condition.GrantExpireAtLte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "grant_expire_at <=",
			Val: condition.GrantExpireAtLte,
		})
	}

	if !string_helper.IsEmpty(condition.ContactName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "contact_name",
			Val: condition.ContactName,
		})
	}
	if !string_helper.IsEmpty(condition.Alias) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "alias",
			Val: condition.Alias,
		})
	}
	if condition.DeptIds != nil && len(condition.DeptIds) > 0 {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dept_id in",
			Val: condition.DeptIds,
		})
	}
	if !string_helper.IsEmpty(condition.LiveAtLte) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "live_at <=",
			Val: condition.LiveAtLte,
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
	if condition.ProcessSwitch != nil && *condition.ProcessSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "process_switch",
			Val: *condition.ProcessSwitch,
		})
	}
	if condition.DllSwitch != nil && *condition.DllSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "dll_switch",
			Val: *condition.DllSwitch,
		})
	}
	if condition.TamperSwitch != nil && *condition.TamperSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "tamper_switch",
			Val: *condition.TamperSwitch,
		})
	}
	if condition.BlackmailSwitch != nil && *condition.BlackmailSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "blackmail_switch",
			Val: *condition.BlackmailSwitch,
		})
	}
	if condition.PortSwitch != nil && *condition.PortSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "port_switch",
			Val: *condition.PortSwitch,
		})
	}
	if condition.PlugSwitch != nil && *condition.PlugSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "plug_switch",
			Val: *condition.PlugSwitch,
		})
	}
	if condition.UsbSwitch != nil && *condition.UsbSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "usb_switch",
			Val: *condition.UsbSwitch,
		})
	}
	if condition.OutreachAlarmSwitch != nil && *condition.OutreachAlarmSwitch != enums.AllDefaultInt {
		wheres = append(wheres, model_comm.Conditions{
			Key: "outreach_alarm_switch",
			Val: *condition.OutreachAlarmSwitch,
		})
	}
	return wheres
}

// ColumnSimplify 字段: 最简单的
func (m *baseAgent) ColumnSimplify() string {
	return "org_id, id, dept_id"
}

// ColumnGrant 字段: 授权想
func (m *baseAgent) ColumnGrant() string {
	column := "org_id, id,"
	column += "grant_product_title,grant_config,grant_order_no,grant_expire_at,grant_status"
	return column
}

// ColumnBasic 字段: 基础
func (m *baseAgent) ColumnBasic(needId, grant bool) string {
	column := ""
	if needId {
		column = "org_id, id,"
	}
	if grant {
		column += "grant_product_title,grant_config, grant_order_no,grant_expire_at,grant_status,"
	}
	column = column +
		"platform_type, os_type, host_name, mac_id, agent_ip, net_ip, os_name, mem_desc, " +
		"hd_desc, cpu_desc, network_card, agent_start_at, os_start_at, ver," +
		"is_online, online_at, live_at, is_safe, dept_id,user_id," +
		"agent_type,alias,contact_name,contact_tel, position, remark, created_at"
	return column
}

// ColumnConfig 字段:全部配置相关
func (m *baseAgent) ColumnConfig(needId, grant bool) string {
	column := ""
	if needId {
		column = "org_id, id, host_name, agent_ip,"
	}
	if grant {
		column += "grant_product_title,grant_config, grant_order_no,grant_expire_at,grant_status,"
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
		"platform_type, os_type, ver, process_switch, process_mode, dll_switch, dll_mode, tamper_switch, tamper_mode, " +
		"blackmail_switch, blackmail_mode, port_switch,plug_switch, usb_switch, usb_mode, outreach_alarm_switch"
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
	data.OfflineNumber = data.TotalNumber - data.OnlineNumber

	return data
}
