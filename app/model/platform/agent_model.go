package model

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 创建 CreateAgentFactory

func CreateAgentFactory() *AgentModel {
	return &AgentModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type AgentModel struct {
	model.BaseModel
	Uuid                 string `json:"uuid"`
	HostName             string `gorm:"column:host_name" json:"host_name"`
	MacId                string `gorm:"column:mac_id" json:"mac_id"`
	Ip                   string `json:"ip"`
	ClientIp             string `gorm:"column:client_ip" json:"client_ip"`
	OsName               string `gorm:"column:os_name" json:"os_name"`
	PlatformType         string `gorm:"column:platform_type" json:"platform_type"`
	OsType               string `gorm:"column:os_type" json:"os_type"`
	MemSize              string `gorm:"column:mem_size" json:"mem_size"`
	HdSize               string `gorm:"column:hd_size" json:"hd_size"`
	CpuSize              string `gorm:"column:cpu_size" json:"cpu_size"`
	AgentStartAt         string `gorm:"column:agent_start_at" json:"agent_start_at"`
	OsStartAt            string `gorm:"column:os_start_at" json:"os_start_at"`
	Ver                  string `json:"ver"`
	Auth                 string `json:"auth"`
	IsOnline             int    `gorm:"column:os_start_at" json:"is_online"`
	OnlineStartAt        string `gorm:"column:online_start_at" json:"online_start_at"`
	LiveAt               string `gorm:"column:live_at" json:"live_at"`
	UpgradeVer           string `gorm:"column:upgrade_ver" json:"upgrade_ver"`
	AdminId              int    `gorm:"column:admin_id" json:"admin_id"`
	DeptId               int    `gorm:"column:dept_id" json:"dept_id"`
	AgentTypeId          int    `gorm:"column:agent_type_id" json:"agent_type_id"`
	AliasName            string `gorm:"column:alias_name" json:"alias_name"`
	ContactName          string `gorm:"column:contact_name" json:"contact_name"`
	ContactTel           string `gorm:"column:contact_tel" json:"contact_tel"`
	Remark               string `gorm:"column:remark" json:"remark"`
	ProcessSwitch        int    `gorm:"column:process_switch" json:"process_switch"`
	ProcessMode          int    `gorm:"column:process_mode" json:"process_mode"`
	DllSwitch            int    `gorm:"column:dll_switch" json:"dll_switch"`
	DllMode              int    `gorm:"column:dll_mode" json:"dll_mode"`
	TamperSwitch         int    `gorm:"column:tamper_switch" json:"tamper_switch"`
	TamperMode           int    `gorm:"column:tamper_mode" json:"tamper_mode"`
	BlackmailSwitch      int    `gorm:"column:blackmail_switch" json:"blackmail_switch"`
	BlackmailMode        int    `gorm:"column:blackmail_mode" json:"blackmail_mode"`
	PortSwitch           int    `gorm:"column:port_switch" json:"port_switch"`
	PlugSwitch           int    `gorm:"column:plug_switch" json:"plug_switch"`
	PlugPortConfig       string `gorm:"column:plug_port_config" json:"plug_port_config"`
	PlugLoginConfig      string `gorm:"column:plug_login_config" json:"plug_login_config"`
	UsbSwitch            int    `gorm:"column:usb_switch" json:"usb_switch"`
	UsbMode              int    `gorm:"column:usb_mode" json:"usb_mode"`
	SyslogDnsDwitch      int    `gorm:"column:syslog_dns_switch" json:"syslog_dns_switch"`
	SyslogOuterSwitch    int    `gorm:"column:syslog_outer_switch" json:"syslog_outer_switch"`
	SyslogSnnerSwitch    int    `gorm:"column:syslog_inner_switch" json:"syslog_inner_switch"`
	SyslogProcessSwitch  int    `gorm:"column:syslog_process_switch" json:"syslog_process_switch"`
	Syslog_loginSwitch   int    `gorm:"column:syslog_login_switch" json:"syslog_login_switch"`
	OutreachAlarm_switch int    `gorm:"column:outreach_alarm_switch" json:"outreach_alarm_switch"`
	SelfProtectionSwitch int    `gorm:"column:self_protection_switch" json:"self_protection_switch"`
	KernelUpgradeSwitch  int    `gorm:"column:kernel_upgrade_switch" json:"kernel_upgrade_switch"`
	CommonTime           int    `gorm:"column:common_time" json:"common_time"`
	FastTime             int    `gorm:"column:fast_time" json:"fast_time"`
	AutoWhiteSwitch      int    `gorm:"column:auto_white_switch" json:"auto_white_switch"`
	AlarmAiSwitch        int    `gorm:"column:alarm_ai_switch" json:"alarm_ai_switch"`
	LoginSolveSwitch     int    `gorm:"column:login_solve_switch" json:"login_solve_switch"`
	LoginSolveConfig     string `gorm:"column:login_solve_config" json:"login_solve_config"`
	HardwareSwitch       string `gorm:"column:hardware_switch" json:"hardware_switch"`
	HardwareTime         string `gorm:"column:hardware_time" json:"hardware_time"`
	ProcessFlushAt       string `gorm:"column:process_flush_at" json:"process_flush_at"`
	DllFlushAt           string `gorm:"column:dll_flush_at" json:"dll_flush_at"`
}

// 表名

func (m *AgentModel) TableName() string {
	return "a_agents"
}

func (m *AgentModel) AddOne(userName, pass, userIp string) bool {
	sql := "INSERT  INTO a_agents(uuid, id, host_name, mac_id, ip, client_ip, os_name, platform_type, os_type, " +
		"mem_size, hd_size, cpu, agent_start_at, os_start_at, ver, auth, is_online, online_start_at, live_at, " +
		"is_safe, upgrade_ver, admin_id, dept_id, agent_type_id, alias_name, contact_name, contact_tel, remark, " +
		"process_switch, process_mode, dll_switch, dll_mode, tamper_switch, tamper_mode, blackmail_switch," +
		" blackmail_mode, port_switch, plug_switch, usb_switch, usb_mode, syslog_dns_switch, syslog_outer_switch, " +
		"syslog_inner_switch, syslog_process_switch, syslog_login_switch, outreach_alarm_switch, " +
		"self_protection_switch, kernel_upgrade_switch, common_time, fast_time, auto_white_switch, alarm_ai_switch, " +
		"login_solve_switch, hardware_switch, hardware_time, plug_port_config, plug_login_config, login_solve_config, " +
		"process_flush_at, dll_flush_at, created_at, updated_at) "
	result := m.Exec(sql, userName, pass, userIp, userName)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
