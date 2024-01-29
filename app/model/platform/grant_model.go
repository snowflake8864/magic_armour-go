package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 授权文件表

func CreateGrantFactory() *GrantModel {
	return &GrantModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type GrantModel struct {
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	Code       string `json:"code"`
	Type       int    `json:"type"`
	ExpireAt   string `json:"expire_at"`
	Status     int    `json:"status"`
	License    string `json:"license"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *GrantModel) TableName() string {
	return "p_grants"
}

func (m *GrantModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			code VARCHAR(128)  NOT NULL DEFAULT '' COMMENT '授权码',
			type TINYINT(1) NOT NULL DEFAULT '0' COMMENT '类型 1 测试 2正式',
			expire_at DATE DEFAULT NULL COMMENT '过期时间',
			status TINYINT(1) NOT NULL DEFAULT '1' COMMENT '状态 1 无效 2 有效',
			license TEXT  COMMENT '授权信息',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NULL,
            PRIMARY KEY (id),
			UNIQUE KEY unique_index (code)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='授权表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
