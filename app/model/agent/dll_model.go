package modelAgent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 模块表

func CreateDllFactory() *DllModel {
	return &DllModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type DllModel struct {
	model.BaseModel
	AgentId         string `gorm:"column:AgentId" json:"agent_id"`
	DllHash         string `gorm:"column:DllHash" json:"dll_hash"`
	Dir             string `json:"dir"`
	Status          int    `json:"status"`
	SyncStatusAt    string `gorm:"column:SyncStatusAt" json:"sync_status_at"`
	GlobalStatus    int    `gorm:"column:GlobalStatus" json:"global_status"`
	SelfStatus      int    `gorm:"column:SelfStatus" json:"self_status"`
	ThreatenStatus  int    `gorm:"column:ThreatenStatus" json:"threaten_status"`
	AutoWhiteSwitch int    `gorm:"column:AutoWhiteSwitch" json:"auto_white_switch"`
	Introduce       string `json:"introduce"`
	Copyright       string `json:"copyright"`
	Attribute       string `json:"attribute"`
	CreatedAt       string `gorm:"column:CreatedAt" json:"created_at"`
}

// 表名

func (m *DllModel) TableName() string {
	return "a_dlls"
}

func (m *DllModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '设备id',
			dll_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '模块:hash',
			dir varchar(512)  NOT NULL DEFAULT '' COMMENT '模块:路径',
			status smallint NOT NULL DEFAULT '0' COMMENT '模块状态(文件中宏定义)',
			sync_status_at datetime NULL COMMENT '同步状态时间(黑白名单下发agent用)',
			global_status tinyint NOT NULL DEFAULT '0' COMMENT '在全局名单中(0未设置 1白名单 2黑名单)',
			self_status tinyint NOT NULL DEFAULT '0' COMMENT '本身状态(0未设置 1白名单 2黑名单)',
			threaten_status tinyint NOT NULL DEFAULT '0' COMMENT '在威胁情报中心中(0未设置 1白名单 2黑名单)',
			auto_white_switch tinyint NOT NULL DEFAULT '0' COMMENT '自动匹配白名单开关(冗余的agent_config.auto_white_switch)',
			introduce varchar(256)  NOT NULL DEFAULT '' COMMENT '模块:描述',
			copyright varchar(256)  NOT NULL DEFAULT '' COMMENT '模块:版权',
			signature varchar(512)  NOT NULL DEFAULT '' COMMENT '模块:签名',
			attribute text  NOT NULL COMMENT '模块:属性',
			created_at datetime NULL,
			INDEX index_dir (dir) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(agent_id, dll_hash) 
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
