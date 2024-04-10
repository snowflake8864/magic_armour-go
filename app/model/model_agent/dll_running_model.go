package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 模块运行中表

func CreateDllRunningFactory() *DllRunningModel {
	return &DllRunningModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type DllRunningModel struct {
	model.BaseModel
	Id              string `json:"id" gorm:"primaryKey"`
	OrgId           int64  `json:"org_id"`
	AgentId         string `json:"agent_id"`
	DllHash         string `json:"dll_hash"`
	Dir             string `json:"dir"`
	Status          int    `json:"status"`
	SyncStatusAt    string `json:"sync_status_at"`
	GlobalStatus    int    `json:"global_status"`
	SelfStatus      int    `json:"self_status"`
	ThreatenStatus  int    `json:"threaten_status"`
	AutoWhiteSwitch int    `json:"auto_white_switch"`
	Introduce       string `json:"introduce"`
	Copyright       string `json:"copyright"`
	Attribute       string `json:"attribute"`
	CreatedAt       string `json:"created_at"`
}

// 表名

func (m *DllRunningModel) TableName() string {
	return "a_dll_running"
}

func (m *DllRunningModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			dll_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '模块:hash',
			proc_id int NOT NULL DEFAULT '0'  COMMENT '进程:id',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			dll_dir varchar(512)  NOT NULL DEFAULT '' COMMENT '模块:路径',
            proc_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '进程:hash',
            proc_dir varchar(512)  NOT NULL DEFAULT '' COMMENT '进程:路径',
			created_at datetime NULL,
			INDEX index_dll_hash (dll_hash) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id, dll_hash, proc_id) 
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
