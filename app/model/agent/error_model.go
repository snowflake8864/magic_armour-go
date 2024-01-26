package modelAgent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 终端错误表

func CreateErrorFactory() *ErrorModel {
	return &ErrorModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ErrorModel struct {
	model.BaseModel
	Id           string `json:"id" gorm:"primaryKey"`
	AgentId      string `json:"agent_id"`
	OriginalName string `json:"original_name"`
	StoragePath  string `json:"storage_path"`
	Remark       string `json:"remark"`
	CreatedAt    string `json:"created_at"`
}

// 表名

func (m *ErrorModel) TableName() string {
	return "a_errors"
}

func (m *ErrorModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(	
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5(created_at+agent_id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			original_name varchar(128)  NOT NULL DEFAULT '' COMMENT '原文件名',
			storage_path varchar(512)  NOT NULL DEFAULT '' COMMENT '存储相对路径(项目下)',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '备注',
			created_at DATETIME NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(id) 
		DISTRIBUTED  BY HASH(id) BUCKETS AUTO 
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
