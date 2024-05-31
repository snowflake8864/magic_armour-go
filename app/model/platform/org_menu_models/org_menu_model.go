package org_menu_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// CreateOrgMenuFactory 企业用的菜单表
func CreateOrgMenuFactory() *OrgMenuModel {
	return &OrgMenuModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OrgMenuModel struct {
	baseMenu
	model.BaseModel
	Id        int64  `gorm:"primaryKey" json:"id"`
	ParentId  int64  `json:"parent_id"`
	Position  int    `json:"position"`
	Type      int    `json:"type"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Sort      int    `json:"sort"`
	IsShow    int    `json:"is_show"`
	IsRoot    int    `json:"is_root"`
	UsageVer  int    `json:"usage_ver"`
	Remark    string `json:"remark"`
	Api       string `json:"api"`
	Level     int    `json:"level"`
	CreatedAt string `json:"created_at"`
}

// 表名

func (m *OrgMenuModel) TableName() string {
	return "p_org_menus"
}

func (m *OrgMenuModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			parent_id int(11) NOT NULL DEFAULT '0' COMMENT '父级id',
			position tinyint(1) NOT NULL DEFAULT '1' COMMENT '方位(1左,2右,3顶部)',
			type tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型(1菜单,2按钮,3链接)',
			name varchar(128)  NOT NULL DEFAULT '' COMMENT '标题',
			code varchar(128)  NOT NULL DEFAULT '' COMMENT '页面标识',
			sort int(11) NOT NULL DEFAULT '1' COMMENT '排序',
			is_show tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否显示(0否,1是),前端判断用',
			is_root tinyint(1) NOT NULL DEFAULT '0' COMMENT '仅超管可用(0否,1是),后端判断用',
			usage_ver tinyint(1) NOT NULL DEFAULT '0' COMMENT '使用版本，二进制表示(1离线版本,2在线版本)',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '备注',
			api text  NOT NULL COMMENT '服务端对应接口(,)隔开',
			level smallint(4) NOT NULL DEFAULT '1' COMMENT '级别(从1开始)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			KEY index_parent (parent_id),
			KEY index_level (level)
		)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) InitData() bool {
	_ = m.TruncateTable()

	sql := `INSERT INTO ` + m.TableName()
	sql += ` (id, parent_id, position, type, name, code, sort, is_show, is_root, usage_ver, remark, api, level, created_at)`
	sql += ` VALUES `
	sql += `
			(1, 1, 1, 1, 'name', 'code', 1, 0, 0, 3,  'remark', 'api', 1, '2024-01-30 00:00:00');
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) InsertItem(s *OrgMenuModel) bool {
	result := m.Create(s)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Table(m.TableName()).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) DeleteById(id int64) bool {
	result := m.Where("id=?", id).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) GetItemById(id int64) (item *OrgMenuModel, err error) {
	result := m.Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return nil, result.Error
	}
}

func (m *OrgMenuModel) GetItemByParent(id int64) (item *OrgMenuModel, err error) {
	result := m.Where("parent_id=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return nil, result.Error
	}
}

func (m *OrgMenuModel) GetItemsByIds(ids []int64) (items []*OrgMenuModel, err error) {
	result := m.Where("id IN ?", ids).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// GetListByConditions 获取符合条件的全部记录
func (m *OrgMenuModel) GetListByConditions(condition SearchConditions) (items []*OrgMenuModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *OrgMenuModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*OrgMenuModel) {
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
