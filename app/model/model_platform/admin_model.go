package model_platform

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
	LoginPhoneCheck int    `json:"login_phone_check"`
	Pw              string `json:"pw"`
	PwErrorNum      int    `json:"pw_error_num"`
	LastLoginAt     string `json:"last_login_at"`
	IsDefault       int    `json:"is_default"`
	IsRoot          int    `json:"is_root"`
	IsVisible       int    `json:"is_visible"`
	CreateId        int    `json:"create_id"`
	CreateName      string `json:"create_name"`
	CreatedAt       string `json:"created_at"`
}

// 表名
func (m *AdminModel) TableName() string {
	return "p_admins"
}

func (m *AdminModel) Columns() []string {
	return []string{
		"id",
		"username",
		"role_id",
		"email",
		"phone",
		"login_phone_check",
		"pw_error_num",
		"last_login_at",
		"is_default",
		"is_root",
		"is_visible",
		"create_id",
		"create_name",
		"created_at",
	}
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
            login_phone_check tinyint(1) NOT NULL DEFAULT '0' COMMENT '手机登录验证(0否,1是)',
			pw varchar(128)  NOT NULL DEFAULT '' COMMENT '密码',
			pw_error_num smallint(4) NOT NULL DEFAULT '0' COMMENT '密码错误次数',
			last_login_at datetime NOT NULL COMMENT '最后登录时间',
			is_default tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认(0否,1是)',
			is_root tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否超管(0否,1是)',
			is_visible tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否可见(0否,1是)',
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
	list := []AdminModel{
		{
			Id:              1,
			Username:        "root",
			RoleId:          0,
			Email:           "",
			Phone:           "",
			LoginPhoneCheck: 0,
			Pw:              "6d030b0eaac5a116f605dc23c5a970e3",
			PwErrorNum:      0,
			LastLoginAt:     curTime,
			IsDefault:       1,
			IsRoot:          1,
			IsVisible:       1,
			CreateId:        0,
			CreateName:      "",
			CreatedAt:       curTime,
		},
	}
	result := m.Create(&list)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *AdminModel) GetItemById(id int) (*AdminModel, error) {
	result := m.Select(m.Columns()).Where("id=?", id).First(m)
	if result.Error == nil {
		return m, nil
	} else {
		return nil, result.Error
	}
}

func (m *AdminModel) GetItemByUsername(username string) (*AdminModel, error) {
	result := m.Select(m.Columns()).Where("username=?", username).First(m)
	if result.Error == nil {
		return m, nil
	} else {
		return nil, result.Error
	}
}

func (m *AdminModel) GetItemByUsernameAndPw(username, pw string) (*AdminModel, error) {
	result := m.Select(m.Columns()).Where("username=?", username).Where("pw=?", pw).First(m)
	if result.Error == nil {
		return m, nil
	} else {
		return nil, result.Error
	}
}

func (m *AdminModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Table(m.TableName()).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}
