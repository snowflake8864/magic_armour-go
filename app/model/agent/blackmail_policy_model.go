package modelAgent

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
	Id         string `json:"id" gorm:"primaryKey"`
	AgentId    string `json:"agent_id"`
	Type       int    `json:"type"`
	TypeVal    string `json:"type_val"`
	GroupId    string `json:"group_id"`
	Processes  string `json:"processes"`
	Remark     string `json:"remark"`
	OpenSwitch int    `json:"open_switch"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
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
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
            created_at datetime NOT NULL,
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			type tinyint NOT NULL DEFAULT '1' COMMENT '策略类型(1后缀 2路径)',
			type_val varchar(512)  NOT NULL DEFAULT '' COMMENT '策略类型值',
			group_id varchar(64)  NOT NULL DEFAULT '' COMMENT '类型组(blackmail_group.id)',
			processes string NOT NULL DEFAULT '' COMMENT '信任进程json',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '备注',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			updated_at datetime NOT NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(agent_id, create_id, created_at) 
		DISTRIBUTED  BY HASH(agent_id) BUCKETS AUTO 
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
