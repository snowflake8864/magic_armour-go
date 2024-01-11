package modelAlarm

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
	AgentId      string `gorm:"column:agent_id" json:"agent_id"`
	Type         int    `json:"type"`
	CreatedAt    string `gorm:"column:created_at" json:"created_at"`
	Level        int    `json:"level"`
	Id           int64  `gorm:"primaryKey" json:"id"`
	AgentAt      string `gorm:"column:agent_at" json:"agent_at"`
	HandleStatus int    `gorm:"column:handle_status" json:"handle_status"`
	HandleRemark string `gorm:"column:handle_remark" json:"handle_remark"`
	HandleUser   string `gorm:"column:handle_user" json:"handle_user"`
	HandleAt     string `gorm:"column:handle_at" json:"handle_at"`
	ProcHash     string `gorm:"column:proc_hash" json:"proc_hash"`
	ProcDir      string `gorm:"column:proc_dir" json:"proc_dir"`
	ProcParam    string `gorm:"column:proc_param" json:"proc_param"`
}

// 表名

func (m *AlarmProcessModel) TableName() string {
	return "a_alarm_process"
}

func (m *AlarmProcessModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id INT NOT NULL DEFAULT '0' COMMENT '终端(agent.id)',
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
			INDEX index_created_at (created_at) USING INVERTED,
			INDEX index_proc_hash (proc_hash) USING INVERTED,
			INDEX index_proc_dir (proc_dir) USING INVERTED
		)
		DUPLICATE  KEY(agent_id, type, created_at, handle_status) 
		PARTITION BY RANGE(created_at)() 
		DISTRIBUTED  BY HASH(agent_id) BUCKETS AUTO 
		PROPERTIES (
			'replication_allocation' = 'tag.location.default: 1',
			'dynamic_partition.enable' = 'true',
			"dynamic_partition.prefix" = "p",
			'dynamic_partition.time_unit' = 'MONTH',
			'dynamic_partition.start' = '-12',
			 "dynamic_partition.end" = "2",
			'dynamic_partition.create_history_partition' = 'true',
			'dynamic_partition.history_partition_num' = '-1'
		);
		`
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
