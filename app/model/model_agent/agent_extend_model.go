package model_agent

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 终端扩展表

func CreateAgentExtendFactory() *AgentExtendModel {
	return &AgentExtendModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type AgentExtendModel struct {
	model.BaseModel
	AgentId          string `json:"agent_id" gorm:"primaryKey"`
	PlugPortConfig   string `json:"plug_port_config"`
	PlugLoginConfig  string `json:"plug_login_config"`
	LoginSolveConfig string `json:"login_solve_config"`
	ProcessFlushAt   string `json:"process_flush_at"`
	DllFlushAt       string `json:"dll_flush_at"`
}

// 表名

func (m *AgentExtendModel) TableName() string {
	return "a_agents"
}

func (m *AgentExtendModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '设备id',
			plug_port_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵:端口虚开策略(json)',
			plug_login_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵:口令破解策略(json)',
			login_solve_config varchar(256)  NOT NULL DEFAULT '' COMMENT '口令破解:策略(json)',
			process_flush_at datetime NULL COMMENT '进程刷新时间',
			dll_flush_at datetime NULL COMMENT '模块刷新时间'
		)
		UNIQUE KEY(agent_id)
		DISTRIBUTED BY HASH(agent_id) BUCKETS AUTO
		PROPERTIES (
			'enable_unique_key_merge_on_write' = 'true',
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

// InsertItem 插入一条
func (m *AgentExtendModel) InsertItem(s *AgentExtendModel) bool {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(s)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateById 通过id修改
func (m *AgentExtendModel) UpdateById(id string, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("agent_id=?", id).Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateByIds 通过id集合修改
func (m *AgentExtendModel) UpdateByIds(ids []string, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("agent_id IN ?", ids).Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateByConditions 通过条件修改
func (m *AgentExtendModel) UpdateByConditions(condition map[string]interface{}, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where(condition).Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteById 通过id删除
func (m *AgentExtendModel) DeleteById(id string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("agent_id=?", id).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByIds 通过id集合删除
func (m *AgentExtendModel) DeleteByIds(ids []string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("agent_id IN ?", ids).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *AgentExtendModel) GetItemById(id, column string) (item *AgentExtendModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Select(column).Where("agent_id=?", id).First(&item)
	})
	result := m.Raw(m.TrimSql(sql)).First(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *AgentExtendModel) GetItemsByIds(ids []string, column string) (items []*AgentExtendModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Select(column).Where("agent_id IN ?", ids).Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}
