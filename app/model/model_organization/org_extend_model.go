package model_organization

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 企业扩展表

func CreateOrgExtendFactory() *OrgExtendModel {
	return &OrgExtendModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OrgExtendModel struct {
	model.BaseModel
	Id                int64  `json:"id" gorm:"primaryKey"`
	Name              string `json:"name"`
	Type              int    `json:"type"`
	UserId            int64  `json:"user_id"`
	UserName          string `json:"user_name"`
	UserEmail         string `json:"user_email"`
	UserPhone         string `json:"user_phone"`
	Address           string `json:"address"`
	Icon              string `json:"icon"`
	Status            int    `json:"status"`
	IsClean           int    `json:"is_clean"`
	IsDel             int    `json:"is_del"`
	Remark            string `json:"remark"`
	LoginPhoneCheck   int    `json:"login_phone_check"`
	HuanshiSwitch     int    `json:"huanshi_switch"`
	TrusteeshipSwitch int    `json:"trusteeship_switch"`
	Source            int    `json:"source"`
	CreateId          int64  `json:"create_id"`
	CreateName        string `json:"create_name"`
	CreatedAt         string `json:"created_at"`
}

// 表名
func (m *OrgExtendModel) TableName() string {
	return "o_org_extend"
}

func (m *OrgExtendModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id INT(11) UNSIGNED NOT NULL COMMENT '同(a_organizations.id)',
			plug_port_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵端口虚开策略(json)',
			plug_login_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵口令破解策略(json)',
			login_solve_config varchar(256)  NOT NULL DEFAULT '' COMMENT '口令破解策略(json)',
			alarm_notice_config VARCHAR(512) DEFAULT '' COMMENT '告警通知策略(json)',
			hardware_notice_config VARCHAR(256) DEFAULT '' COMMENT '硬件通知策略(json)',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
        ) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='企业:扩展表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *OrgExtendModel) InsertItem(s *OrgExtendModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateById 通过id修改
func (m *OrgExtendModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Model(m).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *OrgExtendModel) DeleteById(id int64) bool {
	result := m.Where("id=?", id).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *OrgExtendModel) GetItemById(id int64) (item *OrgExtendModel, err error) {
	result := m.Model(m).Where("id=?", id).First(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *OrgExtendModel) GetItemsByIds(ids []int64) (items []*OrgExtendModel, err error) {
	result := m.Model(m).Where("id IN ?", ids).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}
