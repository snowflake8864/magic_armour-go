package agentModel

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 创建 CreateAgentFactory

func CreateAgentFactory() *AgentModel {
	return &AgentModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type AgentModel struct {
	model.BaseModel
	UserName    string `gorm:"column:user_name" json:"user_name"`
	Pass        string `json:"-"`
	Phone       string `json:"phone"`
	RealName    string `gorm:"column:real_name" json:"real_name"`
	Status      int    `json:"status"`
	Token       string `json:"token"`
	LastLoginIp string `gorm:"column:last_login_ip" json:"last_login_ip"`
}

// 表名

func (m *AgentModel) TableName() string {
	return "a_agents"
}

//创建表

func (m *AgentModel) CreateTable() bool {
	sql := "CREATE TABLE IF NOT EXISTS " + m.TableName()
	sql += "( "
	sql += "`uuid` varchar(128) NOT NULL DEFAULT '' COMMENT '设备id',"
	sql += "`id` int NOT NULL COMMENT '自增id',"
	sql += "`id` int NOT NULL COMMENT '自增id',"
	sql += " );"

	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
