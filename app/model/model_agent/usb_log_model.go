package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 外设插拔日志表

func CreateUsbLogFactory() *UsbLogModel {
	return &UsbLogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type UsbLogModel struct {
	model.BaseModel
	AgentId    string `json:"agent_id"`
	UsbHash    string `json:"usb_hash"`
	Name       string `json:"name"`
	Status     int    `json:"status"`
	PlugStatus int    `json:"plug_status"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *UsbLogModel) TableName() string {
	return "a_usb_logs"
}

func (m *UsbLogModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			usb_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '设备:hash',
			created_at datetime NULL,
			name varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:名称',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(冗余)',
			plug_status tinyint NOT NULL DEFAULT '1' COMMENT '插拔状态(1拔出 2插入)'
		)
		DUPLICATE  KEY(agent_id, usb_hash, created_at)
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
