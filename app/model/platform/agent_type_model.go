package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 终端类型表

func CreateAgentTypeFactory() *AgentTypeModel {
	return &AgentTypeModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type AgentTypeModel struct {
	model.BaseModel
	Id         int    `gorm:"primaryKey" json:"id"`
	Title      string `json:"title"`
	CreateId   int    `gorm:"column:CreateId" json:"create_id"`
	CreateName string `gorm:"column:CreateName"json:"create_name"`
	CreatedAt  string `gorm:"column:created_at" json:"created_at"`
}

// 表名

func (m *AgentTypeModel) TableName() string {
	return "a_agent_type"
}

func (m *AgentTypeModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int NOT NULL COMMENT '自增id',
			title varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
			create_id int(128) NOT NULL DEFAULT '' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NULL
		)
		UNIQUE KEY(id) 
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
