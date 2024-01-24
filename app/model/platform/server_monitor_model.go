package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 短信表

func CreateServerMonitorFactory() *ServerMonitorModel {
	return &ServerMonitorModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体
`monitor_day` Date COMMENT '监控日期(天)',
`server_code` String DEFAULT '' COMMENT '服务器编号',
`cpu_usage` String  DEFAULT '0.00' COMMENT 'cpu:用量%',
`mem_usage` String  DEFAULT '0.00' COMMENT '内存:用量%',
`hd_usage` String  DEFAULT '0.00' COMMENT '硬盘:用量%',
type ServerMonitorModel struct {
	model.BaseModel
	CreatedAt string `gorm:"column:CreatedAt" json:"created_at"`
	CreatedAt string `gorm:"column:CreatedAt" json:"server_code"`
	Phone     string `json:"phone"`
	Msg       string `json:"msg"`
	Code      string `json:"code"`
	Resp      string `json:"resp"`
	Status    int    `json:"status"`
}

// 表名

func (m *ServerMonitorModel) TableName() string {
	return "p_sms_logs"
}

func (m *ServerMonitorModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			phone varchar(16)  NOT NULL DEFAULT '' COMMENT '业务模块',
			created_at datetime NULL,
			msg varchar(256)  NOT NULL DEFAULT '' COMMENT '路径',
			code varchar(256)  NOT NULL DEFAULT '' COMMENT '路径',
			resp varchar(512)  NOT NULL DEFAULT '' COMMENT '路径',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(1成功 0失败)'
		)
		DUPLICATE  KEY(phone, created_at) 
		DISTRIBUTED BY HASH(phone) BUCKETS 1 
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
