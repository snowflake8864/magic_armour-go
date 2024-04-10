package model_agent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 表: 动态封堵

func CreatePlugDynamicFactory() *PlugDynamicModel {
	return &PlugDynamicModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type PlugDynamicModel struct {
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	AgentId    string `json:"agent_id"`
	PlugType   int    `json:"plug_type"`
	Source     int    `json:"source"`
	Ip         string `json:"ip"`
	Direction  int    `json:"direction"`
	Duration   int    `json:"duration"`
	StartAt    string `json:"start_at"`
	EndAt      string `json:"end_at"`
	OpenSwitch int    `json:"open_switch"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *PlugDynamicModel) TableName() string {
	return "a_plug_dynamic"
}

func (m *PlugDynamicModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			created_at datetime NULL,
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			plug_type smallint DEFAULT '0' COMMENT '封堵类型(1临时 2永久)',
			source tinyint NOT NULL DEFAULT '0' COMMENT '封堵来源(1端口虚开)',
			ip varchar(128)  NOT NULL DEFAULT '' COMMENT '封堵ip',
			direction tinyint NOT NULL DEFAULT '0' COMMENT '方向(1IN 2OUT),二进制表示',
			duration int NOT NULL DEFAULT '0' COMMENT '封堵时长(分钟),0为永久',
			start_at datetime DEFAULT NULL COMMENT '封堵开始时间',
			end_at datetime DEFAULT NULL COMMENT '封堵结束时间',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			updated_at datetime NOT NULL,
			INDEX index_id (id) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id, created_at) 
		DISTRIBUTED BY HASH(org_id) BUCKETS AUTO
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
