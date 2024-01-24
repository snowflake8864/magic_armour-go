package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 自定义表格表

func CreateCustomTableFactory() *CustomTableModel {
	return &CustomTableModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type CustomTableModel struct {
	model.BaseModel
	CreateId  int    `gorm:"column:CreateId" json:"create_id"`
	Route     string `json:"route"`
	Cols      string `json:"Cols"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}

// 表名

func (m *CustomTableModel) TableName() string {
	return "p_custom_table"
}

func (m *CustomTableModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			route varchar(128) NOT NULL DEFAULT '' COMMENT '路由',
			cols text  NULL COMMENT '用户自定义列(json)',
			created_at datetime NULL
		)
		UNIQUE KEY(create_id, route) 
		DISTRIBUTED BY HASH(create_id) BUCKETS 1 
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
