package model_agent

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 进程表

func CreateProcessFactory() *ProcessModel {
	return &ProcessModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ProcessModel struct {
	model.BaseModel
	Id              string `json:"id" gorm:"primaryKey"`
	OrgId           int64  `json:"org_id"`
	AgentId         string `json:"agent_id"`
	ProcHash        string `json:"proc_hash"`
	Dir             string `json:"dir"`
	Status          int    `json:"status"`
	SyncStatusAt    string `json:"sync_status_at" gorm:"default:NULL"`
	GlobalStatus    int    `json:"global_status"`
	SelfStatus      int    `json:"self_status"`
	ThreatenStatus  int    `json:"threaten_status"`
	AutoWhiteSwitch int    `json:"auto_white_switch"`
	Introduce       string `json:"introduce"`
	Copyright       string `json:"copyright"`
	Signature       string `json:"signature"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
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
			sync_status_at datetime NULL COMMENT '同步状态时间(黑白名单下发agent用)',
			global_status tinyint NOT NULL DEFAULT '0' COMMENT '在全局名单中(0未设置 1白名单 2黑名单)',
			self_status tinyint NOT NULL DEFAULT '0' COMMENT '本身状态(0未设置 1白名单 2黑名单)',
			threaten_status tinyint NOT NULL DEFAULT '0' COMMENT '在威胁情报中心中(0未设置 1白名单 2黑名单)',
			auto_white_switch tinyint NOT NULL DEFAULT '0' COMMENT '冗余的agent.auto_white_switch',
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
	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// UpdateByConditions 通过条件修改
func (m *ProcessModel) UpdateByConditions(condition map[string]interface{}, updates map[string]interface{}) bool {
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
func (m *ProcessModel) DeleteByConditions(condition map[string]interface{}) bool {
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

// GetItemByOrgAndAgentAndHash 通过企业id\终端id\hash获取一条记录
func (m *ProcessModel) GetItemByOrgAndAgentAndHash(orgId int64, agentId, hash string) (item *ProcessModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND agent_id=? AND proc_hash=?", orgId, agentId, hash).Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetListByConditions 获取符合条件的全部记录
func (m *ProcessModel) GetListByConditions(condition map[string]interface{}) (items []*ProcessModel, err error) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *ProcessModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*ProcessModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
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

// Count 通过条件统计
func (m *ProcessModel) Count(condition map[string]interface{}) int64 {
	var total int64
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Select("id")
		if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
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
func (m *ProcessModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*ProcessModel) {
	total = m.Count(condition)
	if total == 0 {
		return total, nil
	}
	items, err := m.GetItemsByConditions(condition, offset, limit, order)
	if err == nil {
		return total, items
	} else {
		return 0, nil
	}
}
