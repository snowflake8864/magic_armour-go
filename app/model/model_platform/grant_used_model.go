package model_platform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 授权使用记录表

func CreateGrantUsedFactory() *GrantUsedModel {
	return &GrantUsedModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type GrantUsedModel struct {
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	Code       string `json:"code"`
	Type       int    `json:"type"`
	ExpireAt   string `json:"expire_at"`
	Status     int    `json:"status"`
	License    string `json:"license"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *GrantUsedModel) TableName() string {
	return "p_grant_used"
}

func (m *GrantUsedModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			grant_id int(11) NOT NULL DEFAULT '0' COMMENT '授权(grant.id)',
			module varchar(16)  NOT NULL DEFAULT '' COMMENT '模块(grant_module.module)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NULL,
            PRIMARY KEY (id),
			UNIQUE KEY unique_index (grant_id,module,agent_id),
			KEY index_agent_id (agent_id)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='授权使用记录表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
