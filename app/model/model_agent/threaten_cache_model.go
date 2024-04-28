package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 威胁情报缓存表

func CreateThreatenCacheFactory() *ThreatenCacheModel {
	return &ThreatenCacheModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ThreatenCacheModel struct {
	model.BaseModel
	Id        string `json:"id" gorm:"primaryKey"`
	Level     int    `json:"level"`
	Source    int    `json:"source"`
	Dir       string `json:"dir"`
	ParentDir string `json:"parent_dir"`
	Os        string `json:"os"`
	Attribute string `json:"attribute"`
	DataFrom  int    `json:"data_from"`
	CheckAt   string `json:"check_at"`
	CreatedAt string `json:"created_at"`
}

// 表名

func (m *ThreatenCacheModel) TableName() string {
	return "a_threaten_cache"
}

func (m *ThreatenCacheModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'hash',
			level tinyint NOT NULL DEFAULT '0' COMMENT '类型(参考文件中定义)',
			dir varchar(512)  NOT NULL DEFAULT '' COMMENT '文件路径',
			parent_dir varchar(512)  DEFAULT '' COMMENT '上级路径',
			os varchar(128)  DEFAULT '' COMMENT '操作系统',
			attribute string  COMMENT '文件属性',
			data_from tinyint NOT NULL DEFAULT '1' COMMENT '数据来源(1进程 2模块)',
			source tinyint NOT NULL DEFAULT '1' COMMENT '来源',
			check_at datetime DEFAULT NULL COMMENT '验证时间',
			created_at datetime NOT NULL,
			INDEX index_level (level) USING INVERTED,
			INDEX index_dir (dir) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(id) 
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
