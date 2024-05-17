package plug_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 表: 动态封堵

func CreatePlugDynamicFactory() *PlugDynamicModel {
	return &PlugDynamicModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type PlugDynamicModel struct {
	basePlug
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	AgentId    string `json:"agent_id"`
	PlugType   int    `json:"plug_type"`
	Source     int    `json:"source"`
	Ip         string `json:"ip"`
	Direction  int    `json:"direction"`
	Duration   int    `json:"duration"`
	StartAt    string `json:"start_at"`
	EndAt      string `json:"end_at"`
	OpenSwitch int    `json:"open_switch"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *PlugDynamicModel) TableName() string {
	return "a_plug_dynamic"
}

func (m *PlugDynamicModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			created_at datetime NULL,
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			plug_type smallint DEFAULT '0' COMMENT '封堵类型(1临时 2永久)',
			source tinyint NOT NULL DEFAULT '0' COMMENT '封堵来源(1端口虚开)',
			ip varchar(128)  NOT NULL DEFAULT '' COMMENT '封堵ip',
			direction tinyint NOT NULL DEFAULT '0' COMMENT '方向(1IN 2OUT),二进制表示',
			duration int NOT NULL DEFAULT '0' COMMENT '封堵时长(分钟),0为永久',
			start_at datetime DEFAULT NULL COMMENT '封堵开始时间',
			end_at datetime DEFAULT NULL COMMENT '封堵结束时间',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			updated_at datetime NOT NULL,
			INDEX index_id (id) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id, created_at) 
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
func (m *PlugDynamicModel) InsertItem(s *PlugDynamicModel) bool {

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

// UpdateByConditions 通过条件修改
func (m *PlugDynamicModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
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
func (m *PlugDynamicModel) DeleteByConditions(condition SearchConditions) bool {
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

// GetListByConditions 获取符合条件的全部记录
func (m *PlugDynamicModel) GetListByConditions(condition SearchConditions) (items []*PlugDynamicModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *PlugDynamicModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*PlugDynamicModel) {
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
		return nil
	}
}

// Count 通过条件统计
func (m *PlugDynamicModel) Count(condition SearchConditions) int64 {
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
func (m *PlugDynamicModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*PlugDynamicModel) {
	total = m.Count(condition)
	if total == 0 {
		return total, nil
	}
	items = m.GetItemsByConditions(condition, offset, limit, order)
	if len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}

// GetValidByOrgAndAgentAndIp 通过企业id\终端id\hash获取有效记录
func (m *PlugDynamicModel) GetValidByOrgAndAgentAndIp(orgId int64, agentId, ip string) (item *PlugDynamicModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND agent_id=? AND ip=? AND open_switch=?", orgId, agentId, ip, enums.SwitchOn).Take(&item)
	})

	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil && item.OrgId > 0 {
		return item
	} else {
		return nil
	}
}
