package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 访问日志表

func CreateAccessLogFactory() *AccessLogModel {
	return &AccessLogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type AccessLogModel struct {
	model.BaseModel
	AgentId   string `json:"agent_id"`
	Path      string `json:"path"`
	NetIp     string `json:"net_ip"`
	Input     string `json:"input"`
	Status    int    `json:"status"`
	Resp      string `json:"resp"`
	CreatedAt string `json:"created_at"`
}

// 表名

func (m *AccessLogModel) TableName() string {
	return "a_access_logs"
}

func (m *AccessLogModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			path varchar(128)  NOT NULL DEFAULT '' COMMENT '路径',
			created_at datetime NULL,
			net_ip varchar(128)  NOT NULL DEFAULT '' COMMENT '通讯ip',
			input text NULL COMMENT '输入参数',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(1成功 0失败)',
			resp string NOT NULL DEFAULT '' COMMENT '返回结果'
			
		)
		DUPLICATE  KEY(agent_id, path, created_at,net_ip)
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
