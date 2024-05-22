package outreach_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 表: 外联全局名单

func CreateOutreachGlobalFactory() *OutreachGlobalModel {
	return &OutreachGlobalModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type OutreachGlobalModel struct {
	baseOutreach
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	Ip         string `json:"ip"`
	Type       int    `json:"type"`
	Title      string `json:"title"`
	IsDefault  int    `json:"is_default"`
	OpenSwitch int    `json:"open_switch"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *OutreachGlobalModel) TableName() string {
	return "a_outreach_global"
}

func (m *OutreachGlobalModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			ip varchar(128) NOT NULL DEFAULT '' COMMENT 'ip+掩码',
			type tinyint NOT NULL DEFAULT '0' COMMENT '类型(1白 2黑)',
			title varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
			is_default tinyint NOT NULL DEFAULT '0' COMMENT '是否默认',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			INDEX index_ip (ip) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id,id) 
		DISTRIBUTED  BY HASH(org_id) BUCKETS AUTO 
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
func (m *OutreachGlobalModel) InsertItem(s *OutreachGlobalModel) bool {

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
func (m *OutreachGlobalModel) InsertBatch(list *[]OutreachGlobalModel) bool {
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

// DeleteByOrg 通过企业id删除
func (m *OutreachGlobalModel) DeleteByOrg(orgId int64) bool {
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

// DeleteByOrgAndAndId 通过企业id、ID删除
func (m *OutreachGlobalModel) DeleteByOrgAndAndId(orgId int64, id string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Where("id=?", id).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetListByConditions 获取符合条件的全部记录
func (m *OutreachGlobalModel) GetListByConditions(condition SearchConditions) (items []*OutreachGlobalModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *OutreachGlobalModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*OutreachGlobalModel) {
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
func (m *OutreachGlobalModel) Count(condition SearchConditions) int64 {
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
func (m *OutreachGlobalModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*OutreachGlobalModel) {
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
