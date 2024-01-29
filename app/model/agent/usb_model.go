package modelAgent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 外设表

func CreateUsbFactory() *UsbModel {
	return &UsbModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type UsbModel struct {
	model.BaseModel
	Id        string `json:"id" gorm:"primaryKey"`
	AgentId   string `json:"agent_id"`
	UsbHash   string `json:"usb_hash"`
	Name      string `json:"name"`
	Intro     string `json:"intro"`
	Type      string `json:"type"`
	Status    int    `json:"status"`
	IsOnline  int    `json:"is_online"`
	Alias     string `json:"alias"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// 表名

func (m *UsbModel) TableName() string {
	return "a_usb"
}

func (m *UsbModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			usb_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '设备:hash',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			name varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:名称',
			intro varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:简介',
			type varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:类型',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(0默认 1信任 2禁止)',
			is_online tinyint NOT NULL DEFAULT '0' COMMENT '是否在线(0离线 1在线)',
			alias varchar(256)  NOT NULL DEFAULT '' COMMENT '别名',
			remark varchar(256)  NOT NULL DEFAULT '' COMMENT '备注',
			created_at datetime NULL,
			updated_at datetime NULL,
			INDEX index_id (id) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(agent_id, usb_hash) 
		DISTRIBUTED BY HASH(agent_id) BUCKETS AUTO
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
