package agent_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 终端扩展表

func CreateExtendFactory() *ExtendModel {
	return &ExtendModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ExtendModel struct {
	baseAgent
	model.BaseModel
	AgentId          string `json:"id" gorm:"primaryKey;column:id"`
	OrgId            int64  `json:"org_id"`
	PlugPortConfig   string `json:"plug_port_config"`
	PlugLoginConfig  string `json:"plug_login_config"`
	LoginSolveConfig string `json:"login_solve_config"`
	ProcessFlushAt   string `json:"process_flush_at"`
	DllFlushAt       string `json:"dll_flush_at"`
}

// 表名

func (m *ExtendModel) TableName() string {
	return "a_agent_extend"
}

func (m *ExtendModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			id varchar(128) NOT NULL DEFAULT '' COMMENT '设备id',
			plug_port_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵:端口虚开策略(json)',
			plug_login_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵:口令破解策略(json)',
			login_solve_config varchar(256)  NOT NULL DEFAULT '' COMMENT '口令破解:策略(json)',
			process_flush_at datetime NULL COMMENT '进程刷新时间',
			dll_flush_at datetime NULL COMMENT '模块刷新时间',
			INDEX index_org_id (org_id) USING INVERTED
		)
		UNIQUE KEY(org_id, id)
		DISTRIBUTED BY HASH(org_id) BUCKETS AUTO
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
func (m *ExtendModel) InsertItem(s *ExtendModel) bool {

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

// UpdateByOrgAndAgent 通过企业id和终端id修改
func (m *ExtendModel) UpdateByOrgAndAgent(orgId int64, agentId string, updates map[string]interface{}) (bool, int64) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND id=?", orgId, agentId).Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// UpdateByOrgAndAgents 通过企业id和终端id集合修改
func (m *ExtendModel) UpdateByOrgAndAgents(orgId int64, agentIds []string, updates map[string]interface{}) (bool, int64) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND id IN ?", orgId, agentIds).Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// UpdateByConditions 通过条件修改
func (m *ExtendModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) (bool, int64) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// DeleteByOrgAndAgent 通过企业id和终端id删除
func (m *ExtendModel) DeleteByOrgAndAgent(orgId int64, agentId string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND id=?", orgId, agentId).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return false
	}
}

// DeleteByOrg 通过企业i删除
func (m *ExtendModel) DeleteByOrg(orgId int64) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return false
	}
}

// GetItemByOrgAndAgent 通过企业id和终端id获取一条记录
func (m *ExtendModel) GetItemByOrgAndAgent(orgId int64, agentId, column string) (item *ExtendModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Select(m.GetColumn(column)).Where("org_id=? AND id=?", orgId, agentId).Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return item
	}
}
