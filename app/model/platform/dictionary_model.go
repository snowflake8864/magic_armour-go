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
	Id     string `gorm:"primaryKey" json:"id"`
	Module string `json:"module"`
	Val    string `json:"val"`
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
			Id varchar(64) NOT NULL DEFAULT '' COMMENT '标识',
			module varchar(32)  NOT NULL DEFAULT '' COMMENT '模块',
			val varchar(1024)  NOT NULL DEFAULT '' COMMENT '值',
			remark varchar(128)   NOT NULL DEFAULT '' COMMENT '备注',
			PRIMARY KEY (id),
			KEY index_module (module)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据字典表';
		`
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
