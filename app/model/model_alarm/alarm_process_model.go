package model_alarm

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 创建 CreateAlarmProcessFactory

func CreateAlarmProcessFactory() *AlarmProcessModel {
	return &AlarmProcessModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAlarm)}}
}

//结构体

type AlarmProcessModel struct {
	model.BaseModel
	Id           int64  `json:"id" gorm:"primaryKey"`
	AgentId      string `json:"agent_id"`
	Type         int    `json:"type"`
	Level        int    `json:"level"`
	AgentAt      string `json:"agent_at"`
	HandleStatus int    `json:"handle_status"`
	HandleRemark string `json:"handle_remark"`
	HandleUser   string `json:"handle_user"`
	HandleAt     string `json:"handle_at"`
	CreatedAt    string `json:"created_at"`
	ProcHash     string `json:"proc_hash"`
	ProcDir      string `json:"proc_dir"`
	ProcParam    string `json:"proc_param"`
}

// 表名

func (m *AlarmProcessModel) TableName() string {
	return "a_alarm_process"
}

func (m *AlarmProcessModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			type SMALLINT NOT NULL DEFAULT '0' COMMENT '类型',
			created_at DATETIME NULL,
			handle_status TINYINT NOT NULL DEFAULT '0' COMMENT '处理状态(0未处理，1已处理)',
			level TINYINT NOT NULL DEFAULT '0' COMMENT '级别',
			id BIGINT NOT NULL DEFAULT '0',
			agent_at DATETIME NULL COMMENT '终端时间',
			handle_remark VARCHAR(512) NOT NULL DEFAULT '' COMMENT '处理备注',
			handle_user VARCHAR(128) NOT NULL DEFAULT '' COMMENT '处理人(admin.username)',
			handle_at DATETIME NULL COMMENT '处理时间',
			proc_hash VARCHAR(64) NOT NULL DEFAULT '' COMMENT '进程:hash',
			proc_dir VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '进程:路径',
			proc_param VARCHAR(2048) NOT NULL DEFAULT '' COMMENT '进程:路径',
			INDEX index_id (id) USING INVERTED,
			INDEX index_proc_hash (proc_hash) USING INVERTED,
			INDEX index_proc_dir (proc_dir) USING INVERTED
		)
		DUPLICATE  KEY(agent_id, type, created_at, handle_status,level)
		DISTRIBUTED  BY HASH(agent_id) BUCKETS AUTO 
		PROPERTIES (
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
