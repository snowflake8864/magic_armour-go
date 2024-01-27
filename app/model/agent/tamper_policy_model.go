package modelAgent

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

//终端:防篡改策略

func CreateTamperPolicyFactory() *TamperPolicyModel {
	return &TamperPolicyModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type TamperPolicyModel struct {
	model.BaseModel
	Id            string `json:"id" gorm:"primaryKey"`
	AgentId       string `json:"agent_id"`
	Type          int    `json:"type"`
	Dir           string `json:"dir"`
	ProtectRw     int    `json:"protect_rw"`
	IncludeSuffix string `json:"include_suffix"`
	ExcludeSuffix string `json:"exclude_suffix"`
	ExcludeFolder string `json:"exclude_folder"`
	IsExtend      int    `json:"is_extend"`
	Processes     string `json:"processes"`
	OpenSwitch    int    `json:"open_switch"`
	CreateId      int64  `json:"create_id"`
	CreateName    string `json:"create_name"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// 表名

func (m *TamperPolicyModel) TableName() string {
	return "a_tamper_policy"
}

func (m *TamperPolicyModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			type tinyint NOT NULL DEFAULT '0' COMMENT '类型(1文件夹 2单文件)',
			dir varchar(512)  NOT NULL DEFAULT '' COMMENT '文件路径',
			protect_rw int NOT NULL DEFAULT '0' COMMENT '保护权限(1读,2写,4删除,8重命名,16新建)',
			include_suffix varchar(256)  NOT NULL DEFAULT '' COMMENT '包含文件类型(|隔开)',
			exclude_suffix varchar(256)  NOT NULL DEFAULT '' COMMENT '排除文件类型(|隔开)',
			exclude_folder varchar(1024)  NOT NULL DEFAULT '' COMMENT '排除文件夹(|隔开)',
			is_extend tinyint NOT NULL DEFAULT '0' COMMENT '是否向下继承',
			processes string  NOT NULL DEFAULT ''  COMMENT '信任进程json',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			updated_at datetime NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(agent_id, id) 
		DISTRIBUTED BY HASH(agent_id) BUCKETS AUTO
		PROPERTIES (
			'enable_unique_key_merge_on_write' = 'true',
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
