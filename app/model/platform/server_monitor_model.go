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
	Id         int64  `gorm:"primaryKey" json:"id"`
	MonitorDay string `gorm:"column:MonitorDay" json:"monitor_day"`
	ServerCode string `gorm:"column:ServerCode" json:"server_code"`
	CpuUsage   string `gorm:"column:CpuUsage" json:"cpu_usage"`
	MemUsage   string `gorm:"column:MemUsage" json:"mem_usage"`
	HdUsage    string `gorm:"column:HdUsage" json:"hd_usage"`
	CreatedAt  string `gorm:"column:CreatedAt" json:"created_at"`
}

// 表名

func (m *ServerMonitorModel) TableName() string {
	return "p_server_monitor"
}

func (m *ServerMonitorModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id bigint(16) unsigned NOT NULL AUTO_INCREMENT,
			monitor_day date NOT NULL COMMENT '监控日期(天)',
			server_code varchar(32)  NOT NULL DEFAULT '' COMMENT '服务器编号',
			cpu_usage varchar(256)  NOT NULL DEFAULT '' COMMENT 'cpu:用量%',
			mem_usage varchar(256)  NOT NULL DEFAULT '' COMMENT '内存:用量%',
			hd_usage varchar(512)  NOT NULL DEFAULT '' COMMENT '硬盘:用量%',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			KEY index_monitor_day (monitor_day)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='服务器:监控';
		`
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
