package model_agent

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

//  终端任务表

func CreateTaskFactory() *TaskModel {
	return &TaskModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type TaskModel struct {
	model.BaseModel
	Id          string `json:"id" gorm:"primaryKey"`
	AgentId     string `json:"agent_id"`
	Type        int    `json:"type"`
	Remark      string `json:"remark"`
	Status      int    `json:"status"`
	NetIp       string `json:"net_ip"`
	CreatedAt   string `json:"created_at"`
	CompletedAt string `json:"completed_at"`
}

// 表名

func (m *TaskModel) TableName() string {
	return "a_tasks"
}

func (m *TaskModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			type smallint  NOT NULL DEFAULT '0' COMMENT '类型',
			created_at datetime NULL,
 			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '备注',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(0进行中, 1已完成)',
			net_ip varchar(128)  NOT NULL DEFAULT '' COMMENT '通讯ip',
			completed_at datetime NULL COMMENT '任务完成时间'
		)
		UNIQUE KEY(agent_id, type, created_at)
		DISTRIBUTED  BY HASH(agent_id) BUCKETS AUTO 
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

func (m *TaskModel) InsertBatch(list *[]TaskModel) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(list)
	})
	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *TaskModel) UpdateByCondition(conditions, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where(conditions).Updates(updates)
	})
	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
