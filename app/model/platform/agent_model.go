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
	UserName    string `gorm:"column:user_name" json:"user_name"`
	Pass        string `json:"-"`
	Phone       string `json:"phone"`
	RealName    string `gorm:"column:real_name" json:"real_name"`
	Status      int    `json:"status"`
	Token       string `json:"token"`
	LastLoginIp string `gorm:"column:last_login_ip" json:"last_login_ip"`
}

// 表名

func (m *AgentModel) TableName() string {
	return "a_agents"
}

//创建表

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
	platform_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '平台类型(win|linxu, 统一定义)',
	os_type smallint NOT NULL DEFAULT '0' COMMENT '系统类别(参考ubuntu|centos等,统一定义)',
	mem_size varchar(16)  NOT NULL DEFAULT '' COMMENT '内存',
	hd_size varchar(16)  NOT NULL DEFAULT '' COMMENT '硬盘',
	cpu_size varchar(16)  NOT NULL DEFAULT '' COMMENT 'cpu',
	agent_start_at datetime NOT NULL COMMENT '终端启动时间',
	os_start_at datetime NOT NULL COMMENT '系统启动时间',
	ver varchar(32)  NOT NULL DEFAULT '' COMMENT '终端版本号',
	auth varchar(128)  NOT NULL DEFAULT '' COMMENT '安全密码，终端上传',
	is_online tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否在线(0不在,1在)',
	online_start_at datetime NOT NULL COMMENT '上线开始时间',
	live_at datetime NOT NULL COMMENT '最近活动时间',
	is_safe tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否安全(0不安全,1安全)',
	upgrade_ver varchar(32)  NOT NULL DEFAULT '' COMMENT '要升级版本号',
	admin_id int NOT NULL DEFAULT '0' COMMENT '管理员编号(admin.id)',
	dept_id int NOT NULL DEFAULT '0' COMMENT '组织架构编号(department.id)',
	agent_type_id int NOT NULL DEFAULT '0' COMMENT '终端类型(agent_type.id)',
	alias_name varchar(64)  NOT NULL DEFAULT '' COMMENT '别名',
	contact_name varchar(64)  NOT NULL DEFAULT '' COMMENT '责任人:姓名',
	contact_tel varchar(64)  NOT NULL DEFAULT '' COMMENT '责任人:电话',
	remark varchar(256)  NOT NULL DEFAULT '' COMMENT '备注',
	process_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '进程:开关',
	process_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '进程:模式(0监控,1保护)',
	dll_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '模块:开关',
	dll_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '模块:模式(0监控,1保护)',
	tamper_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '防篡改:开关',
	tamper_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '防篡改:模式(0监控,1保护)',
	blackmail_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '勒索:开关',
	blackmail_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '勒索:模式(0监控,1保护)',
	port_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '端口虚开:开关',
	plug_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '封堵:开关',
	plug_port_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵:端口虚开策略(json)',
	plug_login_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵:口令破解策略(json)',
	usb_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '外设:开关',
	usb_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '外设:模式(0监控,1保护)',
	syslog_dns_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '通信日志:dns开关',
	syslog_outer_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '通信日志:对外通讯开关',
	syslog_inner_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '通信日志:对内通讯开关',
	syslog_process_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '通信日志:进程启动开关',
	syslog_login_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '通信日志:登录日志开关',
	outreach_alarm_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '外联告警:开关',
	self_protection_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '自我保护:开关',
	kernel_upgrade_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '内核升级:开关',
	common_time int(11) NOT NULL DEFAULT '60' COMMENT '命令时间间隔:普通(秒)',
	fast_time int(11) NOT NULL DEFAULT '5' COMMENT '命令时间间隔:快速(秒)',
	auto_white_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '自动匹配白名单:开关',
	alarm_ai_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '告警智能处理:开关',
	login_solve_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '口令破解:开关',
	login_solve_config varchar(256)  NOT NULL DEFAULT '' COMMENT '口令破解:策略(json)',
	hardware_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '上传硬件资源:开关',
	hardware_time int(11) NOT NULL DEFAULT '300' COMMENT '上传硬件资源:快速(秒)',
	process_flush_at datetime NULL COMMENT '进程刷新时间',
	dll_flush_at datetime NULL COMMENT '模块刷新时间',
	created_at datetime NULL,
    updated_at datetime NULL,
	INDEX index_id (id) COMMENT 'id',
	INDEX index_ip (ip) COMMENT 'ip',
	INDEX index_admin_id (admin_id) COMMENT 'admin_id',
	INDEX index_dept_id (dept_id) COMMENT 'dept_id',
	INDEX index_os_name (os_name) COMMENT 'os_name',
	INDEX index_platform_type (platform_type) COMMENT 'platform_type'
)
UNIQUE KEY(uuid)
DISTRIBUTED BY HASH(uuid) BUCKETS 1
PROPERTIES (
	"enable_unique_key_merge_on_write" = "true",
	"replication_allocation" = "tag.location.default: 1"
);
`

	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
