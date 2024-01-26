package modelAgent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 威胁情报缓存表

func CreateThreatenPushFactory() *ThreatenPushModel {
	return &ThreatenPushModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ThreatenPushModel struct {
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	Level      int    `json:"level"`
	Dir        string `json:"dir"`
	ParentDir  string `json:"parent_dir"`
	Os         string `json:"os"`
	Introduce  string `json:"introduce"`
	Copyright  string `json:"copyright"`
	Attribute  string `json:"attribute"`
	DataFrom   int    `json:"data_from"`
	PushStatus int    `json:"push_status"`
	PushAt     string `json:"push_at"`
	CreateId   int    `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *ThreatenPushModel) TableName() string {
	return "a_threaten_pushes"
}

func (m *ThreatenPushModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'hash',
			level tinyint NOT NULL DEFAULT '0' COMMENT '类型(参考文件中定义)',
			dir varchar(512)  NOT NULL DEFAULT '' COMMENT '文件路径',
			parent_dir varchar(512)  DEFAULT '' COMMENT '上级路径',
			os varchar(128)  DEFAULT '' COMMENT '操作系统',
			introduce string  COMMENT '描述',
			copyright string  COMMENT '版权',
			attribute string  COMMENT '文件属性',
			data_from tinyint NOT NULL DEFAULT '1' COMMENT '数据来源(1进程 2模块)',
			push_status tinyint NOT NULL DEFAULT '0' COMMENT '送状态(0未推 1已推)',
			push_at datetime DEFAULT NULL COMMENT '推送时间',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
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
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
