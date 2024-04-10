package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 表: 终端正常业务端口

func CreatePortWorkingFactory() *PortWorkingModel {
	return &PortWorkingModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type PortWorkingModel struct {
	model.BaseModel
	Id        string `json:"id" gorm:"primaryKey"`
	OrgId     int64  `json:"org_id"`
	AgentId   string `json:"agent_id"`
	Type      string `json:"type"`
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	Process   string `json:"process"`
	CreatedAt string `json:"created_at"`
}

// 表名

func (m *PortWorkingModel) TableName() string {
	return "a_port_working"
}

func (m *PortWorkingModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			port varchar(8)  NOT NULL DEFAULT '' COMMENT '端口',
            type varchar(8)  NOT NULL DEFAULT '' COMMENT 'ip类型(ipv4,ipv6)',
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			ip varchar(128)  NOT NULL DEFAULT '' COMMENT 'IP地址',
			process varchar(512)  NOT NULL DEFAULT '' COMMENT '进程',
			created_at datetime NOT NULL,
			INDEX index_id (id) USING INVERTED
		)
		UNIQUE KEY(agent_id, port, type) 
		DISTRIBUTED BY HASH(agent_id) BUCKETS AUTO
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
