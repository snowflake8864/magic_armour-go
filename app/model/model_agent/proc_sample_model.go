package model_agent

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 表: 进程与模块样本

func CreateProcSampleFactory() *ProcSampleModel {
	return &ProcSampleModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ProcSampleModel struct {
	model.BaseModel
	Id           string `json:"id" gorm:"primaryKey"`
	OrgId        int64  `json:"org_id"`
	AgentId      string `json:"agent_id"`
	Dir          string `json:"dir"`
	Introduce    string `json:"introduce"`
	Copyright    string `json:"copyright"`
	Signature    string `json:"signature"`
	DataFrom     int    `json:"data_from"`
	Status       int    `json:"status"`
	OriginalName string `json:"original_name"`
	StoragePath  string `json:"storage_path"`
	CreateId     int64  `json:"create_id"`
	CreateName   string `json:"create_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// 表名

func (m *ProcSampleModel) TableName() string {
	return "a_proc_samples"
}

func (m *ProcSampleModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			dir varchar(8)  NOT NULL DEFAULT '' COMMENT '路径',
			introduce varchar(256)  NOT NULL DEFAULT '' COMMENT '描述',
			copyright varchar(256)  NOT NULL DEFAULT '' COMMENT '版权',
			signature varchar(512)  NOT NULL DEFAULT '' COMMENT '签名',
			data_from tinyint NOT NULL DEFAULT '1' COMMENT '数据来源(1进程 2模块)',
			status tinyint NOT NULL DEFAULT '1' COMMENT '状态(0未提取,1提取中,2已提取)',
			original_name varchar(128)  NOT NULL DEFAULT '' COMMENT '原文件名',
			storage_path varchar(512)  NOT NULL DEFAULT '' COMMENT '存储相对路径(项目下)',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
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

func (m *ProcSampleModel) InsertItem(s *ProcSampleModel) bool {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(s)
	})
	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// DeleteByOrg 通过企业id删除
func (m *ProcSampleModel) DeleteByOrg(orgId int64) bool {
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
func (m *ProcSampleModel) DeleteByOrgAndAndId(orgId int64, id string) bool {
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
