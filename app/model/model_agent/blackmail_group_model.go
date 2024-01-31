package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 表: 勒索防护文件类型组

func CreateBlackmailGroupFactory() *BlackmailGroupModel {
	return &BlackmailGroupModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type BlackmailGroupModel struct {
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	Title      string `json:"title"`
	Suffix     string `json:"suffix"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *BlackmailGroupModel) TableName() string {
	return "a_blackmail_group"
}

func (m *BlackmailGroupModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(	
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			title varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
			suffix varchar(128)  NOT NULL DEFAULT '' COMMENT '后缀,英文逗号隔开',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			INDEX index_title (title) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(id) 
		DISTRIBUTED  BY HASH(id) BUCKETS AUTO 
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
