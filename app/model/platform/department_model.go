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
	Id              int    `gorm:"primaryKey" json:"id"`
	Username        string `json:"username"`
	RoleId          int    `gorm:"column:RoleId" json:"role_id"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	PhoneLoginCheck int    `gorm:"column:PhoneLoginCheck" json:"phone_login_check"`
	Pw              string `json:"pw"`
	PwErrorNum      int    `gorm:"column:PwErrorNum" json:"pw_error_num"`
	LastLoginAt     string `gorm:"column:LastLoginAt" json:"last_login_at"`
	IsDefault       int    `gorm:"column:IsDefault" json:"is_default"`
	IsRoot          int    `gorm:"column:IsRoot" json:"is_root"`
	IsVisible       int    `gorm:"column:IsVisible" json:"is_visible"`
	IsAgent         int    `gorm:"column:IsAgent" json:"is_agent"`
	CreateId        int    `gorm:"column:CreateId" json:"create_id"`
	CreateName      string `gorm:"column:CreateName"json:"create_name"`
	CreatedAt       string `gorm:"column:created_at" json:"created_at"`
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
			username varchar(128) NOT NULL DEFAULT '' COMMENT '用户名',
			role_id int NOT NULL DEFAULT '0' COMMENT '角色(role.id)',
			email varchar(128)  NOT NULL DEFAULT '' COMMENT '邮箱',
            phone varchar(16)  NOT NULL DEFAULT '' COMMENT '手机',
            phone_login_check tinyint NOT NULL DEFAULT '0' COMMENT '手机登录验证(0否,1是)',
			pw varchar(128)  NOT NULL DEFAULT '' COMMENT '密码',
			pw_error_num smallint NOT NULL DEFAULT '0' COMMENT '密码错误次数',
			last_login_at datetime NULL COMMENT '最后登录时间',
			is_default tinyint NOT NULL DEFAULT '0' COMMENT '是否默认(0否,1是)',
			is_root tinyint NOT NULL DEFAULT '0' COMMENT '是否超管(0否,1是)',
			is_visible tinyint NOT NULL DEFAULT '0' COMMENT '是否可见(0否,1是)',
			is_agent tinyint NOT NULL DEFAULT '0' COMMENT '是否看全部终端(0否,1是)',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NULL,
			INDEX index_username (username) USING INVERTED
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
