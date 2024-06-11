package threaten_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/enums/enum_threaten"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 威胁情报缓存表

func CreateCacheFactory() *CacheModel {
	return &CacheModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type CacheModel struct {
	baseThreaten
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

func (m *CacheModel) TableName() string {
	return "a_threaten_caches"
}

func (m *CacheModel) CreateTable() bool {
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

func (m *CacheModel) ColumnList() string {
	return "id,level,dir,parent_dir,os,introduce,copyright,signature,attribute,data_from,source,check_at,created_at"
}

func (m *CacheModel) ColumnLevel() string {
	return "id,level"
}

func (m *CacheModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *CacheModel) InsertItem(s *CacheModel) bool {

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

// UpdateById 通过id修改
func (m *CacheModel) UpdateById(id string, updates map[string]interface{}) bool {
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
func (m *CacheModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) (bool, int64) {
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

// DeleteById 通过id删除
func (m *CacheModel) DeleteById(id string) bool {
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
func (m *CacheModel) GetItemById(id string, column string) (item *CacheModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Select(column).Where("id=?", id).Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return nil
	}
}

// GetLevelById 通过id获取等级
func (m *CacheModel) GetLevelById(id string) (isExist bool, level int) {
	model := m.GetItemById(id, m.ColumnLevel())
	if string_helper.IsEmpty(model.Id) {
		return false, enum_threaten.LevelUnknown
	}
	return true, model.Level
}

// GetItemsByIds 通过id集合获取多条记录
func (m *CacheModel) GetItemsByIds(ids []string, column string) (items []*CacheModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Select(column).Where("id IN ?", ids).Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return nil
	}
}

// GetLevelsByIds 通过id集合获取多条记录, 有映射关系
func (m *CacheModel) GetLevelsByIds(ids []string) map[string]int {
	result := make(map[string]int)

	models := m.GetItemsByIds(ids, m.ColumnLevel())
	if models == nil {
		return nil
	}
	for _, model := range models {
		result[model.Id] = model.Level
	}
	return result
}

// GetMapsByIds 通过id集合获取多条记录, 有映射关系
func (m *CacheModel) GetMapsByIds(ids []string) map[string]*CacheModel {
	result := make(map[string]*CacheModel)

	models := m.GetItemsByIds(ids, m.ColumnList())
	if models == nil {
		return nil
	}
	for _, model := range models {
		result[model.Id] = model
	}
	return result
}

// GetItemsByConditions 通过条件获取多条记录
func (m *CacheModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*CacheModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Select(m.ColumnList())
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
func (m *CacheModel) Count(condition SearchConditions) int64 {
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
func (m *CacheModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*CacheModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
