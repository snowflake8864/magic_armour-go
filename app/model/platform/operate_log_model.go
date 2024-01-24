package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 角色表

func CreateOperateLogFactory() *OperateLogModel {
	return &OperateLogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OperateLogModel struct {
	model.BaseModel
	CreatedAt string `gorm:"column:CreatedAt" json:"created_at"`
	Module    string `json:"module"`
	Path      string `json:"path"`
	AdminId   int    `gorm:"column:AdminId" json:"admin_id"`
	AdminName string `gorm:"column:AdminName" json:"admin_name"`
	AgentId   int    `gorm:"column:AgentId" json:"agent_id"`
	AgentUuid string `gorm:"column:AgentUuid" json:"agent_uuid"`
	AgentName string `gorm:"column:AgentName" json:"agent_name"`
	Method    string `json:"method"`
	Ip        string `json:"ip"`
	Input     string `json:"input"`
	Status    int    `json:"status"`
	Response  string `json:"response"`
}

// 表名

func (m *OperateLogModel) TableName() string {
	return "p_operate_logs"
}

func (m *OperateLogModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			created_at datetime NULL,
			module int NOT NULL DEFAULT '0' COMMENT '业务模块',
			path varchar(128)  NOT NULL DEFAULT '' COMMENT '路径',
			admin_id int NOT NULL DEFAULT '0' COMMENT '管理员(admin.id)',
			admin_name varchar(128)  NOT NULL DEFAULT '' COMMENT '管理员(admin.username)',
			agent_id int NOT NULL DEFAULT '0' COMMENT '终端id(agent.id)',
			agent_uuid varchar(64)  NOT NULL DEFAULT '' COMMENT '终端uuid(agent.uuid)',
			agent_name varchar(128)  NOT NULL DEFAULT '' COMMENT '终端名（agent.host_name）',
			method varchar(8)  NOT NULL DEFAULT '' COMMENT '方法',
			ip varchar(128)  NOT NULL DEFAULT '' COMMENT 'ip',
			input text NULL COMMENT '输入参数',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(1成功 0失败)',
			response text NULL COMMENT '返回结果'
		)
		DUPLICATE  KEY(created_at,module,path, admin_id) 
		DISTRIBUTED BY HASH(module) BUCKETS 1 
		PROPERTIES (
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
