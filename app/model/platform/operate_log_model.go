package modelPlatform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 角色表

func CreateOperateLogFactory() *OperateLogModel {
	return &OperateLogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OperateLogModel struct {
	model.BaseModel
	Id         int64  `gorm:"primaryKey" json:"id"`
	Module     string `json:"module"`
	Path       string `json:"path"`
	AgentId    string `gorm:"column:AgentId" json:"agent_id"`
	Method     string `json:"method"`
	Ip         string `json:"ip"`
	Input      string `json:"input"`
	Status     int    `json:"status"`
	Response   string `json:"response"`
	CreateId   int64  `gorm:"column:CreateId" json:"create_id"`
	CreateName string `gorm:"column:CreateName"json:"create_name"`
	CreatedAt  string `gorm:"column:CreatedAt" json:"created_at"`
}

// 表名

func (m *OperateLogModel) TableName() string {
	return "p_operate_logs"
}

func (m *OperateLogModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(	
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			module int(11) NOT NULL DEFAULT '0' COMMENT '业务模块',
			path varchar(128)  NOT NULL DEFAULT '' COMMENT '路径',
			agent_id varchar(128) NOT NULL DEFAULT ''COMMENT '终端id(agent.id)',
			method varchar(8)  NOT NULL DEFAULT '' COMMENT '方法',
			ip varchar(128)  NOT NULL DEFAULT '' COMMENT '操作ip',
			input text NULL COMMENT '输入参数',
			status tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态(1成功 0失败)',
			response text NULL COMMENT '返回结果',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			KEY index_module (module),
			KEY index_path (path),
			KEY index_create_id (create_id),
			KEY index_agent_id (agent_id)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';
		`
	result := m.Exec(m.TrimCreateTable(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
