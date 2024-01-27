package modelSyslog

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 创建 CreateSyslogDnsFactory

func CreateSyslogDnsFactory() *SyslogDnsModel {
	return &SyslogDnsModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAlarm)}}
}

// 结构体
type SyslogDnsModel struct {
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey" `
	AgentId    string `json:"agent_id"`
	AgentAt    string `json:"agent_at"`
	CreatedAt  string `json:"created_at"`
	ProcHash   string `json:"proc_hash"`
	ProcId     int32  `json:"proc_id"`
	ProcDir    string `json:"proc_dir"`
	DomainName string `json:"domain_name"`
	ResIp      string `json:"res_ip"`
}

// 表名

func (m *SyslogDnsModel) TableName() string {
	return "a_syslog_dns"
}

func (m *SyslogDnsModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			created_at DATETIME NULL,
			id BIGINT NOT NULL DEFAULT '0',
			agent_at DATETIME NULL COMMENT '终端时间',
			proc_hash VARCHAR(64) NOT NULL DEFAULT '' COMMENT '进程:hash',
			proc_id INT NOT NULL DEFAULT '0' COMMENT '进程:ID',
			proc_dir VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '进程:路径',
			domain_name VARCHAR(512) NOT NULL DEFAULT '' COMMENT '域名',
			res_ip VARCHAR(512) NOT NULL DEFAULT '' COMMENT '解析ip',
			INDEX index_id (id) USING INVERTED,
			INDEX index_proc_hash (proc_hash) USING INVERTED,
			INDEX index_domain_name (domain_name) USING INVERTED
		)
		DUPLICATE  KEY(agent_id, created_at) 
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
