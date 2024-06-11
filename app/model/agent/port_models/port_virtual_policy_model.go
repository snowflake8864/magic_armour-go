package port_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 表: 终端端口虚开策略

func CreateVirtualPolicyFactory() *VirtualPolicyModel {
	return &VirtualPolicyModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type VirtualPolicyModel struct {
	basePort
	model.BaseModel
	Id          string `json:"id" gorm:"primaryKey"`
	OrgId       int64  `json:"org_id"`
	AgentId     string `json:"agent_id"`
	Type        string `json:"type"`
	Protocol    string `json:"protocol"`
	SourceIp    string `json:"source_ip"`
	SourcePort  string `json:"source_port"`
	ForwardIp   string `json:"forward_ip"`
	ForwardPort string `json:"forward_port"`
	AlarmLevel  int    `json:"alarm_level"`
	OpenSwitch  int    `json:"open_switch"`
	CreateId    int64  `json:"create_id"`
	CreateName  string `json:"create_name"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 表名

func (m *VirtualPolicyModel) TableName() string {
	return "a_port_virtual_policy"
}

func (m *VirtualPolicyModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			type varchar(8)  NOT NULL DEFAULT '' COMMENT 'ip类型(ipv4,ipv6)',
			protocol varchar(8)  NOT NULL DEFAULT '' COMMENT '传输协议(tcp,udp)',
			source_ip varchar(128)  NOT NULL DEFAULT '' COMMENT '源ip',
			source_port varchar(512)  NOT NULL DEFAULT '' COMMENT '源端口:单端口(80);多端口(80,90);范围(80-90)',
			forward_ip varchar(128)  NOT NULL DEFAULT '' COMMENT '转发ip',
			forward_port varchar(64)  NOT NULL DEFAULT '' COMMENT '转发端口',
			alarm_level tinyint NOT NULL DEFAULT '0' COMMENT '告警级别',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id, id) 
		DISTRIBUTED BY HASH(org_id) BUCKETS AUTO
		PROPERTIES (
			'enable_unique_key_merge_on_write' = 'true',
			'replication_allocation' = 'tag.location.default: 1'
		);
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
