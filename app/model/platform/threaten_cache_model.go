package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 威胁情报缓存表

func CreateThreatenCacheFactory() *ThreatenCacheModel {
	return &ThreatenCacheModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type ThreatenCacheModel struct {
	model.BaseModel
	Id        string `gorm:"primaryKey" json:"id"`
	Level     int    `json:"level"`
	Source    int    `json:"source"`
	Dir       string `json:"dir"`
	ParentDir string `gorm:"column:ParentDir" json:"parent_dir"`
	Os        string `json:"os"`
	Attribute string `json:"attribute"`
	DataFrom  int    `gorm:"column:DataFrom" json:"data_from"`
	CheckAt   string `gorm:"column:CheckAt" json:"check_at"`
	CreatedAt string `gorm:"column:createdAt" json:"created_at"`
}

// 表名

func (m *ThreatenCacheModel) TableName() string {
	return "p_threaten_cache"
}

func (m *ThreatenCacheModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'hash',
			level tinyint NOT NULL DEFAULT '0' COMMENT '类型(参考文件中定义)',
			source tinyint NOT NULL DEFAULT '1' COMMENT '来源(1威胁情报中心, 2人工导入, 3内部生成)',
			dir varchar(512)  NOT NULL DEFAULT '' COMMENT '文件路径',
			parent_dir varchar(512)  DEFAULT '' COMMENT '上级路径',
			os varchar(128)  DEFAULT '' COMMENT '操作系统',
			attribute string  COMMENT '文件属性',
			data_from tinyint NOT NULL DEFAULT '1' COMMENT '数据来源(1进程 2模块)',
			check_at datetime DEFAULT NULL COMMENT '验证时间',
			created_at datetime NOT NULL,
			INDEX index_level (level) USING INVERTED,
			INDEX index_dir (dir) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(id) 
		DISTRIBUTED BY HASH(id) BUCKETS 1 
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
