package blackmail_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

//勒索防护-全局

func CreateBlackmailGlobalFactory() *BlackmailGlobalModel {
	return &BlackmailGlobalModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type BlackmailGlobalModel struct {
	model.BaseModel
	Id           int64  `json:"id" gorm:"primaryKey"`
	OrgId        int64  `json:"org_id"`
	PlatformType int    `json:"platform_type"`
	DeptId       int64  `json:"dept_id"`
	Type         int    `json:"type"`
	TypeVal      string `json:"type_val"`
	GroupId      int64  `json:"group_id"`
	Processes    string `json:"processes"`
	OpenSwitch   int    `json:"open_switch"`
	Remark       string `json:"remark"`
	CreateId     int64  `json:"create_id"`
	CreateName   string `json:"create_name"`
	CreatedAt    string `json:"created_at"`
	UpdateId     int64  `json:"update_id"`
	UpdateName   string `json:"update_name"`
	UpdatedAt    string `json:"updated_at"`
}

// 表名

func (m *BlackmailGlobalModel) TableName() string {
	return "a_blackmail_global"
}

func (m *BlackmailGlobalModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			platform_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '平台类型(统一定义)',
			id BIGINT NOT NULL DEFAULT '0',
			dept_id int NOT NULL DEFAULT '0' COMMENT '组织架构编号(department.id)',
			type tinyint NOT NULL DEFAULT '1' COMMENT '策略类型(1后缀 2路径)',
			type_val varchar(512)  NOT NULL DEFAULT '' COMMENT '策略类型值',
			group_id int  NOT NULL DEFAULT '0' COMMENT '类型组(blackmail_group.id)',
			processes string NOT NULL DEFAULT '' COMMENT '信任进程json',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			remark varchar(256)  NOT NULL DEFAULT '' COMMENT '备注',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			update_id int NOT NULL DEFAULT '0' COMMENT '上一次编辑用户(user.id)',
			update_name varchar(128)  NOT NULL DEFAULT '' COMMENT '上一次编辑用户(user.username)',
			updated_at datetime NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, platform_type, id) 
		DISTRIBUTED BY HASH(org_id) BUCKETS AUTO
		PROPERTIES (
			'enable_unique_key_merge_on_write' = 'true',
			'replication_allocation' = 'tag.location.default: 1'
		);
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *BlackmailGlobalModel) InsertItem(s *BlackmailGlobalModel) bool {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(s)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// DeleteByOrg 通过企业id删除
func (m *BlackmailGlobalModel) DeleteByOrg(orgId int64) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByOrgAndAndId 通过企业id、ID删除
func (m *BlackmailGlobalModel) DeleteByOrgAndAndId(orgId int64, id string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Where("id=?", id).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}
