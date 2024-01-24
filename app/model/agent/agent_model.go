package modelAgent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 终端表

func CreateAgentFactory() *AgentModel {
	return &AgentModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type AgentModel struct {
	model.BaseModel
	Uuid                 string `gorm:"unique" json:"uuid"`
	Id                   int    `gorm:"primaryKey" json:"id"`
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
	SyslogDnSwitch       int    `gorm:"column:syslog_dns_switch" json:"syslog_dns_switch"`
	SyslogOuterSwitch    int    `gorm:"column:syslog_outer_switch" json:"syslog_outer_switch"`
	SyslogInnerSwitch    int    `gorm:"column:syslog_inner_switch" json:"syslog_inner_switch"`
	SyslogProcessSwitch  int    `gorm:"column:syslog_process_switch" json:"syslog_process_switch"`
	SyslogLoginSwitch    int    `gorm:"column:syslog_login_switch" json:"syslog_login_switch"`
	OutreachAlarmSwitch  int    `gorm:"column:outreach_alarm_switch" json:"outreach_alarm_switch"`
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

func (m *AgentModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			uuid varchar(128) NOT NULL DEFAULT '' COMMENT '设备id',
			id int NOT NULL COMMENT '自增id',
			host_name varchar(128)  NOT NULL DEFAULT '' COMMENT '主机名',
			mac_id varchar(255)  NOT NULL DEFAULT '' COMMENT '物理地址',
			ip varchar(255)  NOT NULL DEFAULT '' COMMENT '全部ip',
			client_ip varchar(32)  NOT NULL DEFAULT '' COMMENT '连接ip',
			os_name varchar(128)  NOT NULL DEFAULT '' COMMENT '系统名称',
			platform_type tinyint NOT NULL DEFAULT '0' COMMENT '平台类型(win|linxu, 统一定义)',
			os_type smallint NOT NULL DEFAULT '0' COMMENT '系统类别(参考ubuntu|centos等,统一定义)',
			mem_size varchar(16)  NOT NULL DEFAULT '' COMMENT '内存',
			hd_size varchar(16)  NOT NULL DEFAULT '' COMMENT '硬盘',
			cpu varchar(16)  NOT NULL DEFAULT '' COMMENT 'cpu',
			agent_start_at datetime NOT NULL COMMENT '终端启动时间',
			os_start_at datetime NOT NULL COMMENT '系统启动时间',
			ver varchar(32)  NOT NULL DEFAULT '' COMMENT '终端版本号',
			auth varchar(128)  NOT NULL DEFAULT '' COMMENT '安全密码，终端上传',
			is_online tinyint NOT NULL DEFAULT '0' COMMENT '是否在线(0不在,1在)',
			online_start_at datetime NOT NULL COMMENT '上线开始时间',
			live_at datetime NOT NULL COMMENT '最近活动时间',
			is_safe tinyint NOT NULL DEFAULT '1' COMMENT '是否安全(0不安全,1安全)',
			upgrade_ver varchar(32)  NOT NULL DEFAULT '' COMMENT '要升级版本号',
			admin_id int NOT NULL DEFAULT '0' COMMENT '管理员编号(admin.id)',
			dept_id int NOT NULL DEFAULT '0' COMMENT '组织架构编号(department.id)',
			agent_type_id int NOT NULL DEFAULT '0' COMMENT '终端类型(agent_type.id)',
			alias_name varchar(64)  NOT NULL DEFAULT '' COMMENT '别名',
			contact_name varchar(64)  NOT NULL DEFAULT '' COMMENT '责任人:姓名',
			contact_tel varchar(64)  NOT NULL DEFAULT '' COMMENT '责任人:电话',
			remark varchar(256)  NOT NULL DEFAULT '' COMMENT '备注',
			process_switch tinyint NOT NULL DEFAULT '0' COMMENT '进程:开关',
			process_mode tinyint NOT NULL DEFAULT '0' COMMENT '进程:模式(0监控,1保护)',
			dll_switch tinyint NOT NULL DEFAULT '0' COMMENT '模块:开关',
			dll_mode tinyint NOT NULL DEFAULT '0' COMMENT '模块:模式(0监控,1保护)',
			tamper_switch tinyint NOT NULL DEFAULT '0' COMMENT '防篡改:开关',
			tamper_mode tinyint NOT NULL DEFAULT '0' COMMENT '防篡改:模式(0监控,1保护)',
			blackmail_switch tinyint NOT NULL DEFAULT '0' COMMENT '勒索:开关',
			blackmail_mode tinyint NOT NULL DEFAULT '0' COMMENT '勒索:模式(0监控,1保护)',
			port_switch tinyint NOT NULL DEFAULT '0' COMMENT '端口虚开:开关',
			plug_switch tinyint NOT NULL DEFAULT '0' COMMENT '封堵:开关',
			usb_switch tinyint NOT NULL DEFAULT '0' COMMENT '外设:开关',
			usb_mode tinyint NOT NULL DEFAULT '0' COMMENT '外设:模式(0监控,1保护)',
			syslog_dns_switch tinyint NOT NULL DEFAULT '0' COMMENT '通信日志:dns开关',
			syslog_outer_switch tinyint NOT NULL DEFAULT '0' COMMENT '通信日志:对外通讯开关',
			syslog_inner_switch tinyint NOT NULL DEFAULT '0' COMMENT '通信日志:对内通讯开关',
			syslog_process_switch tinyint NOT NULL DEFAULT '0' COMMENT '通信日志:进程启动开关',
			syslog_login_switch tinyint NOT NULL DEFAULT '0' COMMENT '通信日志:登录日志开关',
			outreach_alarm_switch tinyint NOT NULL DEFAULT '0' COMMENT '外联告警:开关',
			self_protection_switch tinyint NOT NULL DEFAULT '0' COMMENT '自我保护:开关',
			kernel_upgrade_switch tinyint NOT NULL DEFAULT '0' COMMENT '内核升级:开关',
			common_time int NOT NULL DEFAULT '60' COMMENT ' 心跳间隔:普通(秒)',
			fast_time int NOT NULL DEFAULT '5' COMMENT '心跳间隔:快速(秒)',
			auto_white_switch tinyint NOT NULL DEFAULT '0' COMMENT '自动匹配白名单:开关',
			alarm_ai_switch tinyint NOT NULL DEFAULT '0' COMMENT '告警智能处理:开关',
			login_solve_switch tinyint NOT NULL DEFAULT '0' COMMENT '口令破解:开关',
			hardware_switch tinyint NOT NULL DEFAULT '0' COMMENT '上传硬件资源:开关',
			hardware_time int NOT NULL DEFAULT '300' COMMENT '上传硬件资源:快速(秒)',
			plug_port_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵:端口虚开策略(json)',
			plug_login_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵:口令破解策略(json)',
			login_solve_config varchar(256)  NOT NULL DEFAULT '' COMMENT '口令破解:策略(json)',
			process_flush_at datetime NULL COMMENT '进程刷新时间',
			dll_flush_at datetime NULL COMMENT '模块刷新时间',
			created_at datetime NULL,
			updated_at datetime NULL,
			INDEX index_id (id) USING INVERTED,
			INDEX index_ip (ip) USING INVERTED,
			INDEX index_admin_id (admin_id) USING INVERTED ,
			INDEX index_dept_id (dept_id) USING INVERTED ,
			INDEX index_os_name (os_name) USING INVERTED,
			INDEX index_platform_type (platform_type) ,
			INDEX index_created_at (created_at)
		)
		UNIQUE KEY(uuid)
		DISTRIBUTED BY HASH(uuid) BUCKETS 1
		PROPERTIES (
			'enable_unique_key_merge_on_write' = 'true',
			'replication_allocation' = 'tag.location.default: 1'
		);
		`
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
