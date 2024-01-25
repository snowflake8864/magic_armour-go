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
	AgentId    string `gorm:"column:agent_id" json:"agent_id"`
	CreatedAt  string `gorm:"column:CreatedAt" json:"created_at"`
	Id         int64  `gorm:"primaryKey" json:"id"`
	AgentAt    string `gorm:"column:agent_at" json:"agent_at"`
	ProcHash   string `gorm:"column:proc_hash" json:"proc_hash"`
	ProcId     int32  `gorm:"column:proc_id" json:"proc_id"`
	ProcDir    string `gorm:"column:proc_dir" json:"proc_dir"`
	DomainName string `gorm:"column:domain_name" json:"domain_name"`
	ResIp      string `gorm:"column:res_ip" json:"res_ip"`
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
			INDEX index_created_at (created_at) USING INVERTED,
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
