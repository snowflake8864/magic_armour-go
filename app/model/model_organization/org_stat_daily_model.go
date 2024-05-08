package model_organization

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// CreateOrgStatDailyFactory 企业统计: 每天的
func CreateOrgStatDailyFactory() *OrgStatDailyModel {
	return &OrgStatDailyModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type OrgStatDailyModel struct {
	model.BaseModel
	OrgId            int64  `json:"org_id"`
	StatDay          string `json:"stat_day"`
	StatMonth        int    `json:"stat_month"`
	UserTotal        int    `json:"user_total"`
	UserOnline       int    `json:"user_online"`
	UserAdd          int    `json:"user_add"`
	AgentTotal       int    `json:"agent_total"`
	AgentOnline      int    `json:"agent_online"`
	AgentAdd         int    `json:"agent_add"`
	AgentLinuxTotal  int    `json:"agent_linux_total"`
	AgentLinuxOnline int    `json:"agent_linux_online"`
	AgentLinuxAdd    int    `json:"agent_linux_add"`
	AgentWinTotal    int    `json:"agent_win_total"`
	AgentWinOnline   int    `json:"agent_win_online"`
	AgentWinAdd      int    `json:"agent_win_add"`
	CreatedAt        string `json:"created_at"`
}

// 表名
func (m *OrgStatDailyModel) TableName() string {
	return "o_org_stat_daily"
}

func (m *OrgStatDailyModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
            stat_month int(11) NOT NULL COMMENT '统计：日期',
			stat_day dateV2 NOT NULL COMMENT '统计：日期',
			user_total int(11) NOT NULL DEFAULT '0' COMMENT '用户:总数',
			user_online int(11) NOT NULL DEFAULT '0' COMMENT '用户:在线数',
			user_add int(11) NOT NULL DEFAULT '0' COMMENT '用户:新增数',
			agent_total int(11) NOT NULL DEFAULT '0' COMMENT '终端:总数',
			agent_online int(11) NOT NULL DEFAULT '0' COMMENT '终端:在线数',
			agent_add int(11) NOT NULL DEFAULT '0' COMMENT '终端:新增数',
			agent_linux_total int(11) NOT NULL DEFAULT '0' COMMENT '终端:linux-总数',
			agent_linux_online int(11) NOT NULL DEFAULT '0' COMMENT '终端:linux-在线数',
			agent_linux_add int(11) NOT NULL DEFAULT '0' COMMENT '终端:linux-新增数',
			agent_win_total int(11) NOT NULL DEFAULT '0' COMMENT '终端:windows-总数',
			agent_win_online int(11) NOT NULL DEFAULT '0' COMMENT '终端:windows-在线数',
			agent_win_add int(11) NOT NULL DEFAULT '0' COMMENT '终端:windows-新增数',
			created_at datetime NULL
		)
		UNIQUE KEY(org_id,stat_month,stat_day)
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
func (m *OrgStatDailyModel) InsertItem(s *OrgStatDailyModel) bool {

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

// InsertBatch 插入多条
func (m *OrgStatDailyModel) InsertBatch(list *[]OrgStatDailyModel) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(list)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateByConditions 通过条件修改
func (m *OrgStatDailyModel) UpdateByConditions(condition map[string]interface{}, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *OrgStatDailyModel) DeleteByConditions(condition map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *OrgStatDailyModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*OrgStatDailyModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Where(condition)
		if limit > 0 {
			tmp = tmp.Offset(offset).Limit(limit)
		}
		if !helpers.StringIsEmpty(order) {
			tmp = tmp.Order(order)
		} else {
			tmp = tmp.Order("created_at ASC")
		}
		return tmp.Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// CalcByConditions 通过条件获取多条记录
func (m *OrgStatDailyModel) CalcByConditions(condition map[string]interface{}, column string) (item *OrgStatDailyModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		whereSQL, vals, _ := m.WhereBuild(condition)
		tmp := tx.Model(m).Select(column).Where(whereSQL, vals...)
		return tmp.Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return nil, result.Error
	}
}

// CalcByConditionsGroupDay 通过条件获取多条记录且按天分组
func (m *OrgStatDailyModel) CalcByConditionsGroupDay(condition map[string]interface{}, column string) (items []*OrgStatDailyModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		whereSQL, vals, _ := m.WhereBuild(condition)
		tmp := tx.Model(m).Select(column).Where(whereSQL, vals...)
		return tmp.Group("stat_day").Order("stat_day ASC").Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}
