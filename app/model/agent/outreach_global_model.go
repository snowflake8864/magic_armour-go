package modelAgent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 表: 外联全局名单

func CreateOutreachGlobalFactory() *OutreachGlobalModel {
	return &OutreachGlobalModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type OutreachGlobalModel struct {
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	Ip         string `json:"ip"`
	Type       int    `json:"type"`
	Title      string `json:"title"`
	IsDefault  int    `json:"is_default"`
	OpenSwitch int    `json:"open_switch"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *OutreachGlobalModel) TableName() string {
	return "a_outreach_global"
}

func (m *OutreachGlobalModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			ip varchar(128) NOT NULL DEFAULT '' COMMENT 'ip+掩码',
			type tinyint NOT NULL DEFAULT '0' COMMENT '类型(1白 2黑)',
			title varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
			is_default tinyint NOT NULL DEFAULT '0' COMMENT '是否默认',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			INDEX index_ip (ip) USING INVERTED,
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
