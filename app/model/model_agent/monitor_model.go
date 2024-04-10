package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

//  终端监控数据

func CreateMonitorFactory() *MonitorModel {
	return &MonitorModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type MonitorModel struct {
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	AgentId    string `json:"agent_id"`
	MonitorDay string `json:"monitor_day"`
	CpuNumber  string `json:"cpu_number"`
	CpuUsage   string `json:"cpu_usage"`
	CpuTops    string `json:"cpu_tops"`
	MemSize    string `json:"mem_size"`
	MemUsage   string `json:"mem_usage"`
	MemTops    string `json:"mem_tops"`
	HdSize     string `json:"hd_size"`
	HdUsage    string `json:"hd_usage"`
	SelfStat   string `json:"self_stat"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *MonitorModel) TableName() string {
	return "a_monitor"
}

func (m *MonitorModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			monitor_day dateV2 NULL  COMMENT '监控日期(天)',
			created_at datetime NULL,
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			cpu_number varchar(64)  NOT NULL  DEFAULT '' COMMENT 'cpu:核数',
			cpu_usage varchar(64)  NOT NULL  DEFAULT '' COMMENT 'cpu:用量%',
			cpu_tops string  NOT NULL  DEFAULT '' COMMENT 'cpu:占用资源最多进程(json)',
			mem_size varchar(64)  NOT NULL  DEFAULT '' COMMENT '内存:大小',
			mem_usage varchar(64)  NOT NULL  DEFAULT '' COMMENT '内存:用量%',
			mem_tops string NOT NULL  DEFAULT '' COMMENT '内存:占用资源最多进程(json)',
			hd_size varchar(64)  NOT NULL  DEFAULT '' COMMENT '硬盘:大小',
			hd_usage varchar(64)  NOT NULL  DEFAULT '' COMMENT '硬盘:用量%',
			self_stat string  DEFAULT '' COMMENT '本身统计(json)'
		)
		DUPLICATE  KEY(org_id, agent_id, monitor_day, created_at)
		DISTRIBUTED BY HASH(agent_id) BUCKETS AUTO 
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
