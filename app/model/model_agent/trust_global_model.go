package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

//全局信任目录

func CreateTrustGlobalFactory() *TrustGlobalModel {
	return &TrustGlobalModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type TrustGlobalModel struct {
	model.BaseModel
	Id           string `json:"id" gorm:"primaryKey"`
	OrgId        int64  `json:"org_id"`
	PlatformType int    `json:"platform_type"`
	Type         int    `json:"type"`
	Dir          string `json:"dir"`
	Name         string `json:"name"`
	IsExtend     int    `json:"is_extend"`
	OpenSwitch   int    `json:"open_switch"`
	CreateId     int64  `json:"create_id"`
	CreateName   string `json:"create_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// 表名

func (m *TrustGlobalModel) TableName() string {
	return "a_trust_global"
}

func (m *TrustGlobalModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			platform_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '平台类型(统一定义)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			type tinyint NOT NULL DEFAULT '0' COMMENT '类型(1文件夹 2单文件)',
			dir varchar(512)  NOT NULL DEFAULT '' COMMENT '文件路径',
			name varchar(256)  NOT NULL DEFAULT '' COMMENT '名称',
			is_extend tinyint NOT NULL DEFAULT '0' COMMENT '是否向下继承',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			updated_at datetime NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id,platform_type,id) 
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
