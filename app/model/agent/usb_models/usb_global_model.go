package usb_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 外设全局名单表

func CreateGlobalFactory() *GlobalModel {
	return &GlobalModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type GlobalModel struct {
	baseUsb
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	Name       string `json:"name"`
	Intro      string `json:"intro"`
	Type       string `json:"type"`
	Status     int    `json:"status"`
	Alias      string `json:"alias"`
	Remark     string `json:"remark"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
	UpdateId   int64  `json:"update_id"`
	UpdateName string `json:"update_name"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *GlobalModel) TableName() string {
	return "a_usb_global"
}

func (m *GlobalModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT '设备:hash',
			name varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:名称',
			intro varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:简介',
			type varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:类型',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(0默认 1信任 2禁止)',
			alias varchar(256)  NOT NULL DEFAULT '' COMMENT '别名',
			remark varchar(256)  NOT NULL DEFAULT '' COMMENT '备注',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			update_id int NOT NULL DEFAULT '0' COMMENT '上一次编辑用户(user.id)',
			update_name varchar(128)  NOT NULL DEFAULT '' COMMENT '上一次编辑用户(user.username)',
			updated_at datetime NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id,id) 
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

func (m *GlobalModel) InsertItem(s *GlobalModel) bool {

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
func (m *GlobalModel) DeleteByOrg(orgId int64) bool {
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
func (m *GlobalModel) DeleteByOrgAndAndId(orgId int64, id string) bool {
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
