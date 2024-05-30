package agent_monitor_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

//  终端监控数据

func CreateMonitorFactory() *MonitorModel {
	return &MonitorModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type MonitorModel struct {
	baseMonitor
	model.BaseModel
	Id         int64   `json:"id" gorm:"primaryKey"`
	OrgId      int64   `json:"org_id"`
	AgentId    string  `json:"agent_id"`
	MonitorDay string  `json:"monitor_day"`
	CpuNumber  string  `json:"cpu_number"`
	CpuUsage   float64 `json:"cpu_usage"`
	CpuTops    string  `json:"cpu_tops"`
	MemSize    string  `json:"mem_size"`
	MemUsage   float64 `json:"mem_usage"`
	MemTops    string  `json:"mem_tops"`
	HdSize     string  `json:"hd_size"`
	HdUsage    float64 `json:"hd_usage"`
	SelfStat   string  `json:"self_stat"`
	CreatedAt  string  `json:"created_at"`
}

type MonitorTop struct {
	Id       int     `json:"id"`
	Dir      string  `json:"dir"`
	Hash     string  `json:"hash"`
	CpuUsage float64 `json:"cpu_usage"`
	MemSize  string  `json:"mem_size"`
	User     string  `json:"user"`
}
type MonitorSelf struct {
	CpuUsage float64 `json:"cpu_usage"`
	MemSize  string  `json:"mem_size"`
}

// 表名

func (m *MonitorModel) TableName() string {
	return "a_monitors"
}

func (m *MonitorModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			monitor_day dateV2 NULL  COMMENT '监控日期(天)',
			created_at datetime NULL,
			id BIGINT NOT NULL DEFAULT '0',
			cpu_number varchar(64)  NOT NULL  DEFAULT '' COMMENT 'cpu:核数',
			cpu_usage FLOAT  NOT NULL  DEFAULT '0.00' COMMENT 'cpu:用量%',
			cpu_tops string  NOT NULL  DEFAULT '' COMMENT 'cpu:占用资源最多进程(json)',
			mem_size varchar(64)  NOT NULL  DEFAULT '' COMMENT '内存:大小',
			mem_usage FLOAT NOT NULL  DEFAULT '0.00' COMMENT '内存:用量%',
			mem_tops string NOT NULL  DEFAULT '' COMMENT '内存:占用资源最多进程(json)',
			hd_size varchar(64)  NOT NULL  DEFAULT '' COMMENT '硬盘:大小',
			hd_usage FLOAT  NOT NULL  DEFAULT '0.00' COMMENT '硬盘:用量%',
			self_stat string  DEFAULT '' COMMENT '本身统计(json)'
		)
		DUPLICATE  KEY(org_id, agent_id, monitor_day, created_at)
		DISTRIBUTED BY HASH(agent_id) BUCKETS AUTO 
		PROPERTIES (
			'replication_allocation' = 'tag.location.default: 1'
		);
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *MonitorModel) InsertItem(s *MonitorModel) bool {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(s)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// DeleteByOrg 通过企业id删除
func (m *MonitorModel) DeleteByOrg(orgId int64) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByOrgAndAgent 通过企业id和终端id删除
func (m *MonitorModel) DeleteByOrgAndAgent(orgId int64, agentId string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Where("agent_id=?", agentId).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByLteCreateAt 通过(创建时间小于节点)删除
func (m *MonitorModel) DeleteByLteCreateAt(at string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("created_at <= ?", at).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetAllByConditions 通过条件获取全部记录
func (m *MonitorModel) GetAllByConditions(condition SearchConditions, order string) (items []*MonitorModel) {
	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	if !string_helper.IsEmpty(order) {
		tmp = tmp.Order(order)
	} else {
		tmp = tmp.Order("created_at ASC")
	}
	result := tmp.Find(&items)
	if result.Error == nil {
		return items
	} else {
		return nil
	}
}
