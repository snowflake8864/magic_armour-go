package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"time"
)

// 管理员表

func CreateAdminFactory() *AdminModel {
	return &AdminModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type AdminModel struct {
	model.BaseModel
	Id              int64  `json:"id" gorm:"primaryKey"`
	Username        string `json:"username"`
	RoleId          int64  `json:"role_id"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	PhoneLoginCheck int    `json:"phone_login_check"`
	Pw              string `json:"pw"`
	PwErrorNum      int    `json:"pw_error_num"`
	LastLoginAt     string `json:"last_login_at"`
	IsDefault       int    `json:"is_default"`
	IsRoot          int    `json:"is_root"`
	IsVisible       int    `json:"is_visible"`
	IsAgent         int    `json:"is_agent"`
	CreateId        int    `json:"create_id"`
	CreateName      string `json:"create_name"`
	CreatedAt       string `json:"created_at"`
}

// 表名

func (m *AdminModel) TableName() string {
	return "p_admins"
}

func (m *AdminModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			username varchar(128) NOT NULL DEFAULT '' COMMENT '用户名',
			role_id int(11) NOT NULL DEFAULT '0' COMMENT '角色(role.id)',
			email varchar(128)  NOT NULL DEFAULT '' COMMENT '邮箱',
            phone varchar(16)  NOT NULL DEFAULT '' COMMENT '手机',
            phone_login_check tinyint(1) NOT NULL DEFAULT '0' COMMENT '手机登录验证(0否,1是)',
			pw varchar(128)  NOT NULL DEFAULT '' COMMENT '密码',
			pw_error_num smallint(4) NOT NULL DEFAULT '0' COMMENT '密码错误次数',
			last_login_at datetime NOT NULL COMMENT '最后登录时间',
			is_default tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认(0否,1是)',
			is_root tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否超管(0否,1是)',
			is_visible tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否可见(0否,1是)',
			is_agent tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否看全部终端(0否,1是)',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY unique_index (username)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *AdminModel) InitData() bool {
	curTime := time.Now().Format(variable.DateFormat)
	sql := " INSERT INTO " + m.TableName() +
		"(`id`, `username`, `role_id`, `email`, `phone`, `phone_login_check`," +
		" `pw`, " +
		"`pw_error_num`, `last_login_at`, `is_default`, `is_root`, `is_visible`, `is_agent`, " +
		"`create_id`, `create_name`, `created_at`)"
	sql += "VALUES ("
	sql += "1, 'root', 0, '', '', 0, " +
		"'18464611f92f2a687aa0636486d4aa2f49a2546860af4208bc83fc42e7a5c603de52e3c886b290676d26b2714f09deb53f6366c085a7833e5177ce9444bc042d', " +
		"0, '" + curTime + "', 1, 1, 0, 1, " +
		"0, '', '" + curTime + "'"
	sql += "), ("
	sql += "2, 'admin', 0, '', '', 0, " +
		"'2806f9ee82bf9b27357f2fb88259b8858acf6c3b6939dadb27cc3e8bff709f8154d41976ebf64f71cfa9547c624dff812b48a7376f96e38a8b82924fb98c2585', " +
		"0, '" + curTime + "', 1, 0, 1, 1, " +
		"0, '', '" + curTime + "'"
	sql += ");"

	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
