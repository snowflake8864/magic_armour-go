package model_organization

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
	Id         int64  `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	Title      string `json:"title"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *AgentTypeModel) TableName() string {
	return "p_agent_type"
}

func (m *AgentTypeModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			title varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NULL,
            PRIMARY KEY (id),
			KEY index_org (org_id),
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='终端类型表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
