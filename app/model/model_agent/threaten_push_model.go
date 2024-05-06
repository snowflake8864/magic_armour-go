package model_agent

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 威胁情报缓存表

func CreateThreatenPushFactory() *ThreatenPushModel {
	return &ThreatenPushModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ThreatenPushModel struct {
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	Level      int    `json:"level"`
	Dir        string `json:"dir"`
	ParentDir  string `json:"parent_dir"`
	Os         string `json:"os"`
	Introduce  string `json:"introduce"`
	Copyright  string `json:"copyright"`
	Signature  string `json:"signature"`
	Attribute  string `json:"attribute"`
	DataFrom   int    `json:"data_from"`
	PushStatus int    `json:"push_status"`
	PushAt     string `json:"push_at" gorm:"default:NULL"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *ThreatenPushModel) TableName() string {
	return "a_threaten_pushes"
}

func (m *ThreatenPushModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'hash',
			level tinyint NOT NULL DEFAULT '0' COMMENT '类型(参考文件中定义)',
			dir varchar(512)  NOT NULL DEFAULT '' COMMENT '文件路径',
			parent_dir varchar(512)  DEFAULT '' COMMENT '上级路径',
			os varchar(128)  DEFAULT '' COMMENT '操作系统',
			introduce varchar(256)  NOT NULL DEFAULT '' COMMENT '描述',
			copyright varchar(256)  NOT NULL DEFAULT '' COMMENT '版权',
			signature varchar(512)  NOT NULL DEFAULT '' COMMENT '签名',
			attribute string NOT NULL DEFAULT ''  COMMENT '属性',
			data_from tinyint NOT NULL DEFAULT '1' COMMENT '数据来源(1进程 2模块)',
			push_status tinyint NOT NULL DEFAULT '0' COMMENT '送状态(0未推 1已推)',
			push_at datetime DEFAULT NULL COMMENT '推送时间',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
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

func (m *ThreatenPushModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// InsertBatch 插入多条
func (m *ThreatenPushModel) InsertBatch(list *[]ThreatenPushModel) bool {
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

// UpdateByIds 通过id集合修改
func (m *ThreatenPushModel) UpdateByIds(ids []string, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("id IN ?", ids).Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *ThreatenPushModel) GetItemsByIds(ids []string) (items []*ThreatenPushModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("id IN ?", ids).Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *ThreatenPushModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*ThreatenPushModel, err error) {
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
func (m *ThreatenPushModel) Count(condition map[string]interface{}) int64 {
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
func (m *ThreatenPushModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*ThreatenPushModel) {
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
