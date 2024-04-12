package model_organization

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 用户表

func CreateUserFactory() *UserModel {
	return &UserModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type UserModel struct {
	model.BaseModel
	Id          int64  `json:"id" gorm:"primaryKey"`
	Username    string `json:"username"`
	OrgId       int64  `json:"org_id"`
	RoleId      int64  `json:"role_id"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Password    string `json:"password"`
	Realname    string `json:"realname"`
	Position    string `json:"position"`
	Intro       string `json:"intro"`
	PwErrorNum  int    `json:"pw_error_num"`
	LastLoginAt string `json:"last_login_at"`
	IsRoot      int    `json:"is_root"`
	CreateId    int64  `json:"create_id"`
	CreateName  string `json:"create_name"`
	CreatedAt   string `json:"created_at"`
}

// 表名

func (m *UserModel) TableName() string {
	return "o_users"
}

func (m *UserModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			username varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
			org_id int(11) NOT NULL DEFAULT '0' COMMENT '企业编号(o_organizations.id)',
			role_id int(11) NOT NULL DEFAULT '0' COMMENT '角色id(o_role.id)',
			email varchar(64) NOT NULL DEFAULT '' COMMENT '邮箱',
			phone varchar(16) NOT NULL DEFAULT '' COMMENT '手机',
			password varchar(128) NOT NULL DEFAULT '' COMMENT '密码',
			realname varchar(128) NOT NULL DEFAULT '' COMMENT '真实姓名',
			position varchar(128) NOT NULL DEFAULT '' COMMENT '职位',
			intro varchar(512) NOT NULL DEFAULT '' COMMENT '简介',
			pw_error_num smallint(4) NOT NULL DEFAULT '0' COMMENT '密码错误次数',
			last_login_at datetime DEFAULT NULL COMMENT '最后登录时间',
			is_root tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否管理员(0否,1是)',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人id(o_users.id)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY username (username),
			KEY org_id (org_id),
			KEY email (email),
			KEY phone (phone)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='企业:用户';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *UserModel) InsertItem(s *UserModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateById 通过id修改
func (m *UserModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Model(m).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *UserModel) DeleteByConditions(condition map[string]interface{}) bool {
	result := m.Model(m).Where(condition).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *UserModel) GetItemById(id int64) (item *UserModel, err error) {
	result := m.Model(m).Where("id=?", id).First(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *UserModel) GetItemsByIds(ids []int64) (items []*UserModel, err error) {
	result := m.Model(m).Where("id IN ?", ids).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *UserModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*UserModel, err error) {
	tmp := m.Model(m).Where(condition)
	if limit > 0 {
		tmp = tmp.Offset(offset).Limit(limit)
	}
	if !helpers.StringIsEmpty(order) {
		tmp = tmp.Order(order)
	} else {
		tmp = tmp.Order("created_at ASC")
	}
	result := m.Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// Count 通过条件统计
func (m *UserModel) Count(condition map[string]interface{}) int64 {
	var total int64
	result := m.Model(m).Select("id").Where(condition).Count(&total)
	if result.Error == nil {
		return total
	} else {
		return 0
	}
}

// Page 通过条件分页
func (m *UserModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*UserModel) {
	total = m.Count(condition)
	if total == 0 {
		return total, nil
	}
	items, err := m.GetItemsByConditions(condition, offset, limit, order)
	if err == nil {
		return total, items
	} else {
		return 0, nil
	}
}
