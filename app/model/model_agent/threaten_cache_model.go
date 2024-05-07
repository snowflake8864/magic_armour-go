package model_agent

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 威胁情报缓存表

func CreateThreatenCacheFactory() *ThreatenCacheModel {
	return &ThreatenCacheModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ThreatenCacheModel struct {
	model.BaseModel
	Id        string `json:"id" gorm:"primaryKey"`
	Level     int    `json:"level"`
	Source    int    `json:"source"`
	Dir       string `json:"dir"`
	ParentDir string `json:"parent_dir"`
	Os        string `json:"os"`
	Introduce string `json:"introduce"`
	Copyright string `json:"copyright"`
	Signature string `json:"signature"`
	Attribute string `json:"attribute"`
	DataFrom  int    `json:"data_from"`
	CheckAt   string `json:"check_at" gorm:"default:NULL"`
	CheckInfo string `json:"check_info"`
	CreatedAt string `json:"created_at"`
}

// 表名

func (m *ThreatenCacheModel) TableName() string {
	return "a_threaten_caches"
}

func (m *ThreatenCacheModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'hash',
			level tinyint NOT NULL DEFAULT '0' COMMENT '类型(参考文件中定义)',
			dir varchar(1024)  NOT NULL DEFAULT '' COMMENT '文件路径',
			parent_dir varchar(1024)  DEFAULT '' COMMENT '上级路径',
			os varchar(128)  DEFAULT ''  COMMENT '操作系统',
			introduce varchar(256)  NOT NULL DEFAULT '' COMMENT '描述',
			copyright varchar(256)  NOT NULL DEFAULT '' COMMENT '版权',
			signature varchar(512)  NOT NULL DEFAULT '' COMMENT '签名',
			attribute string NOT NULL DEFAULT '' COMMENT '属性',
			data_from tinyint NOT NULL DEFAULT '1' COMMENT '数据来源(1进程 2模块)',
			source tinyint NOT NULL DEFAULT '1' COMMENT '来源',
			check_at datetime DEFAULT NULL COMMENT '验证时间',
			check_info string NOT NULL DEFAULT '' COMMENT '验证信息',
			created_at datetime NOT NULL,
			INDEX index_level (level) USING INVERTED,
			INDEX index_dir (dir) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(id) 
		DISTRIBUTED BY HASH(id) BUCKETS AUTO 
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

func (m *ThreatenCacheModel) getListColumns() string {
	return "id,level,dir,parent_dir,os,introduce,copyright,signature,attribute,data_from,source,check_at,created_at"
}

func (m *ThreatenCacheModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *ThreatenCacheModel) InsertItem(s *ThreatenCacheModel) bool {

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

// UpdateById 通过id修改
func (m *ThreatenCacheModel) UpdateById(id string, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("id=?", id).Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateByConditions 通过条件修改
func (m *ThreatenCacheModel) UpdateByConditions(condition map[string]interface{}, updates map[string]interface{}) bool {
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

// DeleteById 通过id删除
func (m *ThreatenCacheModel) DeleteById(id string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("id=?", id).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *ThreatenCacheModel) GetItemById(id string) (item *ThreatenCacheModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("id=?", id).Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *ThreatenCacheModel) GetItemsByIds(ids []string) (items []*ThreatenCacheModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Select(m.getListColumns()).Where("id IN ?", ids).Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *ThreatenCacheModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*ThreatenCacheModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Select(m.getListColumns())
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
func (m *ThreatenCacheModel) Count(condition map[string]interface{}) int64 {
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
func (m *ThreatenCacheModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*ThreatenCacheModel) {
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
