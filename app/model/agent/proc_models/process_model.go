package proc_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 进程表

func CreateProcessFactory() *ProcessModel {
	return &ProcessModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ProcessModel struct {
	baseProc
	model.BaseModel
	Id                  string `json:"id" gorm:"primaryKey"`
	OrgId               int64  `json:"org_id"`
	AgentId             string `json:"agent_id"`
	ProcHash            string `json:"proc_hash"`
	Dir                 string `json:"dir"`
	Status              int    `json:"status"`
	SyncStatusAt        int64  `json:"sync_status_at"`
	GlobalStatus        int    `json:"global_status"`
	SelfStatus          int    `json:"self_status"`
	ThreatenStatus      int    `json:"threaten_status"`
	ThreatenLevel       int    `json:"threaten_level"` //未使用
	AutoWhiteSwitch     int    `json:"auto_white_switch"`
	GrantThreatenSwitch int    `json:"grant_threaten_switch"`
	Introduce           string `json:"introduce"`
	Copyright           string `json:"copyright"`
	Signature           string `json:"signature"`
	CreatedAt           string `json:"created_at"`
}

// 表名

func (m *ProcessModel) TableName() string {
	return "a_processes"
}

func (m *ProcessModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			proc_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '进程:hash',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			dir varchar(512)  NOT NULL DEFAULT '' COMMENT '进程:路径',
			status smallint NOT NULL DEFAULT '0' COMMENT '进程状态(文件中宏定义)',
			sync_status_at int NOT NULL DEFAULT '0' COMMENT '同步状态时间(黑白名单下发agent用)',
			global_status tinyint NOT NULL DEFAULT '0' COMMENT '在全局名单状态中',
			self_status tinyint NOT NULL DEFAULT '0' COMMENT '本身状态',
			threaten_status tinyint NOT NULL DEFAULT '0' COMMENT '在威胁情报中状态',
			threaten_level tinyint NOT NULL DEFAULT '0' COMMENT '威胁情报等级(冗余threaten_cache.level)',
			auto_white_switch tinyint NOT NULL DEFAULT '0' COMMENT '自动白名单开关(冗余agent.auto_white_switch)',
			grant_threaten_switch tinyint NOT NULL DEFAULT '0' COMMENT '授权威胁情报开关(冗余agent.grant_config)',
			introduce varchar(256)  NOT NULL DEFAULT '' COMMENT '进程:描述',
			copyright varchar(256)  NOT NULL DEFAULT '' COMMENT '进程:版权',
			signature varchar(512)  NOT NULL DEFAULT '' COMMENT '进程:签名',
			created_at datetime NULL,
			INDEX index_dir (dir) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id,agent_id, proc_hash) 
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

func (m *ProcessModel) InsertItem(s *ProcessModel) bool {

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

// UpdateByConditions 通过条件修改
func (m *ProcessModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) (bool, int64) {
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

// DeleteByConditions 通过条件删除
func (m *ProcessModel) DeleteByConditions(condition SearchConditions) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
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

// GetItemByOrgAndAgentAndHash 通过企业id\终端id\hash获取一条记录
func (m *ProcessModel) GetItemByOrgAndAgentAndHash(orgId int64, agentId, hash string) (item *ProcessModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND agent_id=? AND proc_hash=?", orgId, agentId, hash).Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item
	} else {
		return nil
	}
}

// GetListByConditions 获取符合条件的全部记录
func (m *ProcessModel) GetListByConditions(condition SearchConditions) (items []*ProcessModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *ProcessModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*ProcessModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		if limit > 0 {
			tmp = tmp.Offset(offset).Limit(limit)
		}
		if !string_helper.IsEmpty(order) {
			tmp = tmp.Order(order)
		} else {
			tmp = tmp.Order("created_at ASC")
		}
		return tmp.Find(&items)
	})

	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return nil
	}
}

// Count 通过条件统计
func (m *ProcessModel) Count(condition SearchConditions) int64 {
	var total int64
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Select("id")
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Count(&total)
	})
	result := m.Raw(m.TrimSql(sql)).Count(&total)
	if result.Error == nil {
		return total
	} else {
		return 0
	}
}

// Page 通过条件分页
func (m *ProcessModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*ProcessModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}

// CalcByConditionsGroupHash 通过条件获取多条记录且按hash分组
func (m *ProcessModel) CalcByConditionsGroupHash(condition SearchConditions) (items []*RepeatCalc) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Select(m.ColumnRepeatCalc(true))
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Group("proc_hash").Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items
	} else {
		return nil
	}
}

// CalcByConditionsGroupAgent 通过条件获取多条记录且按终端分组
func (m *ProcessModel) CalcByConditionsGroupAgent(condition SearchConditions, column string) (items []*StatusCalc) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Select(column)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Group("agent_id").Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items
	} else {
		return nil
	}
}

// GetListByConditionsDistinctHash 通过条件获取多条记录
func (m *ProcessModel) GetListByConditionsDistinctHash(condition SearchConditions) (items []*ProcessModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		tmp.Select("distinct(proc_hash), dir, agent_id")
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}

		return tmp.Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return nil
	}
}
