package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"time"
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

func (m *OutreachGlobalModel) InitData() bool {
	curTime := time.Now().Format(variable.DateFormat)
	sql := " INSERT INTO " + m.TableName() +
		"(`id`, `ip`, `type`, `title`, `is_default`, `open_switch`, " +
		"`create_id`, `create_name`, `created_at`, `updated_at`)"
	sql += "VALUES "
	sql += "('382d29ad257f55c55c75c9f3a906c99e', 'fe80::1/10', '1', '系统默认', 1, 1, 0, '', '" + curTime + "' , '" + curTime + "'),"
	sql += "('0f39d743e689c10ca2fe6e8f57d3d2cd', '127.0.0.0/8', '1', '系统默认', 1, 1, 0, '', '" + curTime + "' , '" + curTime + "'),"
	sql += "('738a5b3b5ed9abce50d289ee49772ddd', '192.168.0.0/16', '1', '系统默认', 1, 1, 0, '', '" + curTime + "' , '" + curTime + "'),"
	sql += "('bb6febec1d21cbfd72a76d25a1ef6622', '172.16.0.0/12', '1', '系统默认', 1, 1, 0, '', '" + curTime + "' , '" + curTime + "'),"
	sql += "('920069446bf4f6b26db744f88f1106bb', '10.0.0.0/8', '1', '系统默认', 1, 1, 0, '', '" + curTime + "' , '" + curTime + "');"

	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
