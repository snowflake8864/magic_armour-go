package model_organization

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
	Id         int64  `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	Route      string `json:"route"`
	Cols       string `json:"cols"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *CustomTableModel) TableName() string {
	return "p_custom_table"
}

func (m *CustomTableModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			route varchar(64) NOT NULL DEFAULT '' COMMENT '路由',
			cols text DEFAULT NULL COMMENT '用户自定义列(json)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY unique_index (org_id,create_id, route)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='自定义表格表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
