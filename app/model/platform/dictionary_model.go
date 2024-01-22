package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 创建 数据字典表

func CreateDictionaryFactory() *DictionaryModel {
	return &DictionaryModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type DictionaryModel struct {
	model.BaseModel
	Uuid   string `gorm:"primaryKey" json:"uuid"`
	Module string `json:"module"`
	Value  string `json:"value"`
	Remark string `json:"remark"`
}

// 表名

func (m *DictionaryModel) TableName() string {
	return "p_dictionaries"
}

func (m *DictionaryModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			uuid varchar(64) NOT NULL DEFAULT '' COMMENT '标识',
			module varchar(32)  NOT NULL DEFAULT '' COMMENT '模块',
			value varchar(1024)  NOT NULL DEFAULT '' COMMENT '值',
			remark varchar(128)   NOT NULL DEFAULT '' COMMENT '备注'
		)
		UNIQUE KEY(uuid) 
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
