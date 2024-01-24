package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 短信表

func CreateServerMonitorFactory() *ServerMonitorModel {
	return &ServerMonitorModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

// 结构体
type ServerMonitorModel struct {
	model.BaseModel
	MonitorDay string `gorm:"column:MonitorDay" json:"monitor_day"`
	ServerCode string `gorm:"column:ServerCode" json:"server_code"`
	CreatedAt  string `gorm:"column:CreatedAt" json:"created_at"`
	CpuUsage   string `gorm:"column:CpuUsage" json:"cpu_usage"`
	MemUsage   string `gorm:"column:MemUsage" json:"mem_usage"`
	HdUsage    string `gorm:"column:HdUsage" json:"hd_usage"`
}

// 表名

func (m *ServerMonitorModel) TableName() string {
	return "p_server_monitor"
}

func (m *ServerMonitorModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			monitor_day dateV2 NULL COMMENT '监控日期(天)',
			server_code varchar(32)  NOT NULL DEFAULT '' COMMENT '服务器编号',
			created_at datetime NULL,
			cpu_usage varchar(256)  NOT NULL DEFAULT '' COMMENT 'cpu:用量%',
			mem_usage varchar(256)  NOT NULL DEFAULT '' COMMENT '内存:用量%',
			hd_usage varchar(512)  NOT NULL DEFAULT '' COMMENT '硬盘:用量%'
		)
		DUPLICATE  KEY(monitor_day,server_code,created_at) 
		DISTRIBUTED BY HASH(monitor_day) BUCKETS 1 
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
