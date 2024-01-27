package modelAgent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 进程运行中表

func CreateProcessRunningFactory() *ProcessRunningModel {
	return &ProcessRunningModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ProcessRunningModel struct {
	model.BaseModel
	Id        string `json:"id" gorm:"primaryKey"`
	AgentId   string `json:"agent_id"`
	ProcHash  string `json:"proc_hash"`
	ProcId    string `json:"proc_id"`
	ProcDir   string `json:"proc_dir"`
	ProcUser  string `json:"proc_user"`
	DllNumber int    `json:"dll_number"`
	DllStatus int    `json:"dll_status"`
	DllList   string `json:"dll_list"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// 表名

func (m *ProcessRunningModel) TableName() string {
	return "a_process_running"
}

func (m *ProcessRunningModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			proc_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '进程:hash',
			proc_id int NOT NULL DEFAULT '0'  COMMENT '进程:id',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
            proc_dir varchar(512)  NOT NULL DEFAULT '' COMMENT '进程:路径',
            proc_user varchar(128)  NOT NULL DEFAULT '' COMMENT '进程:用户',
            dll_number int  NOT NULL DEFAULT '0' COMMENT '模块:数量',
            dll_status tinyint  NOT NULL DEFAULT '0' COMMENT '模块:状态(0默认 1开始, 2获取中, 3结束)',
            dll_list string  NOT NULL DEFAULT '' COMMENT '模块:列表',
			created_at datetime NULL,
			updated_at datetime NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(agent_id, proc_hash, proc_id) 
		DISTRIBUTED BY HASH(agent_id) BUCKETS AUTO
		PROPERTIES (
			'enable_unique_key_merge_on_write' = 'true',
			'replication_allocation' = 'tag.location.default: 1'
		);
		`
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
