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
	Id         int    `gorm:"primaryKey" json:"id"`
	ParentId   int    `gorm:"column:ParentId" json:"parent_id"`
	Name       string `json:"name"`
	FullCode   string `gorm:"column:FullCode" json:"full_code"`
	FullName   string `gorm:"column:FullName" json:"full_name"`
	Ip         string `json:"ip"`
	Level      int    `json:"level"`
	Sort       int    `json:"sort"`
	CreateId   int    `gorm:"column:CreateId" json:"create_id"`
	CreateName string `gorm:"column:CreateName"json:"create_name"`
	CreatedAt  string `gorm:"column:created_at" json:"created_at"`
}

// 表名

func (m *DepartmentModel) TableName() string {
	return "p_departments"
}

func (m *DepartmentModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int NOT NULL DEFAULT '0' COMMENT '自增id',
			parent_id int NOT NULL DEFAULT '0' COMMENT '父级id',
			name varchar(128)  NOT NULL DEFAULT '' COMMENT '名称',
			full_code varchar(128)  NOT NULL DEFAULT '' COMMENT '完整编号(根,父,本级)',
			full_name varchar(2056)  NOT NULL DEFAULT '' COMMENT '完整名称(json)',
			ip varchar(512)  NOT NULL DEFAULT '' COMMENT 'IP段',
			level smallint NOT NULL DEFAULT '1' COMMENT '级别(从1开始)',
			sort int NOT NULL DEFAULT '1' COMMENT '排序',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NULL,
			INDEX index_parent (parent_id) USING INVERTED,
			INDEX index_full_code (full_code) USING INVERTED
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
