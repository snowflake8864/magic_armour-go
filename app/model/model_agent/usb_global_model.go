package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 外设全局名单表

func CreateUsbGlobalFactory() *UsbGlobalModel {
	return &UsbGlobalModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type UsbGlobalModel struct {
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	Intro      string `json:"intro"`
	Type       string `json:"type"`
	Status     int    `json:"status"`
	IsOnline   int    `json:"is_online"`
	Alias      string `json:"alias"`
	Remark     string `json:"remark"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *UsbGlobalModel) TableName() string {
	return "a_usb_global"
}

func (m *UsbGlobalModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id varchar(64)  NOT NULL DEFAULT '' COMMENT '设备:hash',
			name varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:名称',
			intro varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:简介',
			type varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:类型',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(0默认 1信任 2禁止)',
			is_online tinyint NOT NULL DEFAULT '0' COMMENT '是否在线(0离线 1在线)',
			alias varchar(256)  NOT NULL DEFAULT '' COMMENT '别名',
			remark varchar(256)  NOT NULL DEFAULT '' COMMENT '备注',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			updated_at datetime NULL,
			INDEX index_id (id) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(id) 
		DISTRIBUTED BY HASH(id) BUCKETS AUTO
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
