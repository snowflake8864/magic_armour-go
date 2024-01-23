package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 角色表

func CreateRoleFactory() *RoleModel {
	return &RoleModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type RoleModel struct {
	model.BaseModel
	Id         int    `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	Remark     string `json:"remark"`
	IsDefault  int    `gorm:"column:IsDefault" json:"is_default"`
	CreateId   int    `gorm:"column:CreateId" json:"create_id"`
	CreateName string `gorm:"column:CreateName" json:"create_name"`
	CreatedAt  string `gorm:"column:created_at" json:"created_at"`
}

// 表名

func (m *RoleModel) TableName() string {
	return "p_roles"
}

func (m *RoleModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int NOT NULL DEFAULT '0' COMMENT '自增id',
			name varchar(128) NOT NULL DEFAULT '' COMMENT '名称',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '邮箱',
			is_default tinyint NOT NULL DEFAULT '0' COMMENT '是否默认(0否,1是)',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NULL,
			INDEX index_name (name) USING INVERTED
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
