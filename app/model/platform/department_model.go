package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 组织架构表

func CreateDepartmentFactory() *DepartmentModel {
	return &DepartmentModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type DepartmentModel struct {
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	ParentId   int64  `json:"parent_id"`
	Name       string `json:"name"`
	FullCode   string `json:"full_code"`
	FullName   string `json:"full_name"`
	Ip         string `json:"ip"`
	Level      int    `json:"level"`
	Sort       int    `json:"sort"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *DepartmentModel) TableName() string {
	return "p_departments"
}

func (m *DepartmentModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			parent_id int(11) NOT NULL DEFAULT '0' COMMENT '父级id',
			name varchar(128)  NOT NULL DEFAULT '' COMMENT '名称',
			full_code varchar(128)  NOT NULL DEFAULT '' COMMENT '完整编号(根,父,本级)',
			full_name varchar(2056)  NOT NULL DEFAULT '' COMMENT '完整名称(json)',
			ip varchar(512)  NOT NULL DEFAULT '' COMMENT 'IP段',
			level smallint(4) NOT NULL DEFAULT '1' COMMENT '级别(从1开始)',
			sort int(11) NOT NULL DEFAULT '1' COMMENT '排序',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			KEY index_parent (parent_id),
			KEY index_full_code (full_code)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='组织架构表';
		`
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
