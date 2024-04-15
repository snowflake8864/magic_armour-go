package model_organization

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 组织架构表

func CreateDepartmentFactory() *DepartmentModel {
	return &DepartmentModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type DepartmentModel struct {
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	ParentId   int64  `json:"parent_id"`
	Name       string `json:"name"`
	FullCode   string `json:"full_code"`
	FullName   string `json:"full_name"`
	Ip         string `json:"ip"`
	Level      int    `json:"level"`
	Sort       int    `json:"sort"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名
func (m *DepartmentModel) TableName() string {
	return "o_departments"
}

func (m *DepartmentModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			parent_id int(11) NOT NULL DEFAULT '0' COMMENT '父级id',
			name varchar(128)  NOT NULL DEFAULT '' COMMENT '名称',
			full_code varchar(128)  NOT NULL DEFAULT '' COMMENT '完整编号(根,父,本级)',
			full_name varchar(2056)  NOT NULL DEFAULT '' COMMENT '完整名称(json)',
			ip varchar(512)  NOT NULL DEFAULT '' COMMENT 'IP段',
			level smallint(4) NOT NULL DEFAULT '1' COMMENT '级别(从1开始)',
			sort int(11) NOT NULL DEFAULT '1' COMMENT '排序',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			KEY index_org (org_id),
			KEY index_parent (parent_id),
			KEY index_full_code (full_code)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='企业:组织架构表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *DepartmentModel) InsertItem(s *DepartmentModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateById 通过id修改
func (m *DepartmentModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Model(m).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *DepartmentModel) DeleteByConditions(condition map[string]interface{}) bool {
	result := m.Model(m).Where(condition).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *DepartmentModel) GetItemById(id int64) (item *DepartmentModel, err error) {
	result := m.Model(m).Where("id=?", id).First(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *DepartmentModel) GetItemsByIds(ids []int64) (items []*DepartmentModel, err error) {
	result := m.Model(m).Where("id IN ?", ids).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *DepartmentModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*DepartmentModel, err error) {
	tmp := m.Model(m).Where(condition)
	if limit > 0 {
		tmp = tmp.Offset(offset).Limit(limit)
	}
	if !helpers.StringIsEmpty(order) {
		tmp = tmp.Order(order)
	} else {
		tmp = tmp.Order("created_at ASC")
	}
	result := m.Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}
