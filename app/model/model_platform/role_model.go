package model_platform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"time"
)

// 角色表

func CreateRoleFactory() *RoleModel {
	return &RoleModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type RoleModel struct {
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	Remark     string `json:"remark"`
	Menus      string `json:"menus"`
	IsDefault  int    `json:"is_default"`
	CreateId   int64  `json:"create_id"`
	CreateName string `son:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *RoleModel) TableName() string {
	return "p_roles"
}

func (m *RoleModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			name varchar(128) NOT NULL DEFAULT '' COMMENT '名称',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '邮箱',
			menus text DEFAULT NULL COMMENT '邮箱',
			is_default tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认(0否,1是)',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
            KEY index_name (name)
		)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *RoleModel) InitData() bool {
	curTime := time.Now().Format(variable.DateFormat)
	sql := " INSERT INTO " + m.TableName() +
		"(`id`, `name`, `remark`, `menus`, `is_default`, " +
		"`create_id`, `create_name`, `created_at`)"
	sql += "VALUES "
	sql += "(1, '管理员', '系统默认', '[]', 1, 0, '',  '" + curTime + "');"

	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
