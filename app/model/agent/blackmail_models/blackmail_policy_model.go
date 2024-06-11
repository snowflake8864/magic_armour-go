package blackmail_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 勒索防护策略表

func CreateBlackmailPolicyFactory() *BlackmailPolicyModel {
	return &BlackmailPolicyModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type BlackmailPolicyModel struct {
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	AgentId    string `json:"agent_id"`
	OrgId      int64  `json:"org_id"`
	Type       int    `json:"type"`
	TypeVal    string `json:"type_val"`
	GroupId    int64  `json:"group_id"`
	Processes  string `json:"processes"`
	OpenSwitch int    `json:"open_switch"`
	Remark     string `json:"remark"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
	UpdateId   int64  `json:"update_id"`
	UpdateName string `json:"update_name"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *BlackmailPolicyModel) TableName() string {
	return "a_blackmail_policy"
}

func (m *BlackmailPolicyModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(	
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			id BIGINT NOT NULL DEFAULT '0',
			type tinyint NOT NULL DEFAULT '1' COMMENT '策略类型(1后缀 2路径)',
			type_val varchar(512)  NOT NULL DEFAULT '' COMMENT '策略类型值',
			group_id int  NOT NULL DEFAULT '0' COMMENT '类型组(blackmail_group.id)',
			processes string NOT NULL DEFAULT '' COMMENT '信任进程json',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			remark varchar(256)  NOT NULL DEFAULT '' COMMENT '备注',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			update_id int NOT NULL DEFAULT '0' COMMENT '上一次编辑用户(user.id)',
			update_name varchar(128)  NOT NULL DEFAULT '' COMMENT '上一次编辑用户(user.username)',
			updated_at datetime NOT NULL
		)
		UNIQUE KEY(org_id, agent_id, id) 
		DISTRIBUTED  BY HASH(org_id) BUCKETS AUTO 
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
