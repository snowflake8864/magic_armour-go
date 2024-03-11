package model_platform

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
	Id         int64  `json:"id" gorm:"primaryKey"`
	Path       string `json:"path"`
	Module     string `json:"module"`
	NetIp      string `json:"net_ip"`
	Input      string `json:"input"`
	AgentId    string `json:"agent_id"`
	Method     string `json:"method"`
	Status     int    `json:"status"`
	HttpCode   int    `json:"http_code"`
	DelayTime  int64  `json:"delay_time"`
	Resp       string `json:"resp"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
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
			path varchar(128)  NOT NULL DEFAULT '' COMMENT '路径',
			module varchar(32) NOT NULL DEFAULT '' COMMENT '业务模块',
			net_ip varchar(128) NOT NULL DEFAULT '' COMMENT '通讯ip',
			input text NULL COMMENT '请求参数',
			agent_id varchar(128) NOT NULL DEFAULT ''COMMENT '终端id(agent.id)',
			method varchar(8) NULL COMMENT '请求方法',
			status tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态(1成功 0失败)',
			http_code smallint(4) NOT NULL DEFAULT '0' COMMENT 'http状态码',
			delay_time int NOT NULL DEFAULT '0' COMMENT '延迟时间(ms)',
			resp text NULL COMMENT '返回结果',
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
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
