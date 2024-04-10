package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 表: 终端封堵IP名单库

func CreatePlugIpPolicyFactory() *PlugIpPolicyModel {
	return &PlugIpPolicyModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type PlugIpPolicyModel struct {
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	AgentId    string `json:"agent_id"`
	Ip         string `json:"ip"`
	Type       int    `json:"type"`
	Duration   int    `json:"duration"`
	Remark     string `json:"remark"`
	OpenSwitch int    `json:"open_switch"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *PlugIpPolicyModel) TableName() string {
	return "a_plug_ip_policy"
}

func (m *PlugIpPolicyModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			ip varchar(128)  NOT NULL DEFAULT '' COMMENT 'ip',
			type tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型(1白 2黑)',
			direction tinyint(1) NOT NULL DEFAULT '0' COMMENT '方向(1IN 2OUT),二进制表示',
			remark varchar(256)  NOT NULL DEFAULT '' COMMENT '备注',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			INDEX index_agent_id (agent_id) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id, id) 
		DISTRIBUTED BY HASH(id) BUCKETS AUTO
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
