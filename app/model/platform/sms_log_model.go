package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 短信表

func CreateSmsLogFactory() *SmsLogModel {
	return &SmsLogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type SmsLogModel struct {
	model.BaseModel
	CreatedAt string `gorm:"column:createdAt" json:"created_at"`
	Phone     string `json:"phone"`
	Msg       string `json:"msg"`
	Code      string `json:"code"`
	Resp      string `json:"resp"`
	Status    int    `json:"status"`
}

// 表名

func (m *SmsLogModel) TableName() string {
	return "p_sms_logs"
}

func (m *SmsLogModel) CreateTable() bool {
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
