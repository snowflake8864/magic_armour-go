package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 授权模块表

func CreateGrantModuleFactory() *GrantModuleModel {
	return &GrantModuleModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type GrantModuleModel struct {
	model.BaseModel
	Id      int64  `json:"id" gorm:"primaryKey"`
	GrantId int64  `json:"int64"`
	Module  string `json:"module"`
	Total   int    `json:"status"`
	Used    int    `json:"used"`
}

// 表名

func (m *GrantModuleModel) TableName() string {
	return "p_grant_module"
}

func (m *GrantModuleModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			grant_id int(11) NOT NULL DEFAULT '0' COMMENT '授权码id(grant.id)',
			module varchar(16)  NOT NULL DEFAULT '' COMMENT '相应模块(参考宏)',
			total int(11) NOT NULL DEFAULT '0' COMMENT '授权总数',
			used int(11) NOT NULL DEFAULT '0' COMMENT '已用数量',
            PRIMARY KEY (id),
			UNIQUE KEY unique_index (grant_id, module)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='授权模块表';
		`
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
