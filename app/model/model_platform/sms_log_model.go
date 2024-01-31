package model_platform

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
	Id        int64  `json:"id" gorm:"primaryKey"`
	Phone     string `json:"phone"`
	Msg       string `json:"msg"`
	Code      string `json:"code"`
	Resp      string `json:"resp"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

// 表名

func (m *SmsLogModel) TableName() string {
	return "p_sms_logs"
}

func (m *SmsLogModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			phone varchar(16)  NOT NULL DEFAULT '' COMMENT '业务模块',
			msg varchar(256)  NOT NULL DEFAULT '' COMMENT '内容',
			code varchar(256)  NOT NULL DEFAULT '' COMMENT '验证码',
			resp varchar(512)  NOT NULL DEFAULT '' COMMENT '结果',
			status tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态(1成功 0失败)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			KEY index_phone (phone)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='短信日志表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
