package model_agent

import (
	"gorm.io/gorm"
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
	Id        string `json:"id" gorm:"primaryKey"`
	OrgId     int64  `json:"org_id"`
	AgentId   string `json:"agent_id"`
	Path      string `json:"path"`
	NetIp     string `json:"net_ip"`
	Input     string `json:"input"`
	Method    string `json:"method"`
	Status    int    `json:"status"`
	HttpCode  int    `json:"http_code"`
	DelayTime int64  `json:"delay_time"`
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
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			path varchar(128)  NOT NULL DEFAULT '' COMMENT '请求地址',
			created_at datetime NULL,
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			net_ip varchar(128)  NOT NULL DEFAULT '' COMMENT '通讯ip',
			input text NULL COMMENT '请求参数',
			method varchar(8) NULL COMMENT '请求方法',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(1成功 0失败)',
			http_code smallint NOT NULL DEFAULT '0' COMMENT 'http状态码',
			delay_time int NOT NULL DEFAULT '0' COMMENT '延迟时间(ms)',
			resp string NOT NULL DEFAULT '' COMMENT '返回结果'
		)
		DUPLICATE  KEY(org_id, agent_id, path, created_at)
		DISTRIBUTED  BY HASH(org_id) BUCKETS AUTO 
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

func (m *AccessLogModel) InsertItem(s *AccessLogModel) bool {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(s)
	})
	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
