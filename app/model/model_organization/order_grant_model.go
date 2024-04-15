package model_organization

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 订单表

func CreateOrderGrantFactory() *OrderGrantModel {
	return &OrderGrantModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OrderGrantModel struct {
	model.BaseModel
	OrderNo    string `json:"order_no" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	AgentTotal int    `json:"agent_total"`
	AgentUsed  int    `json:"agent_used"`
}

// 表名
func (m *OrderGrantModel) TableName() string {
	return "o_order_grant"
}

func (m *OrderGrantModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			order_no VARCHAR(128) NOT NULL DEFAULT '' COMMENT '订单编号',
			org_id int(11) NOT NULL DEFAULT '0' COMMENT '企业(o_organization.id)',
			agent_total int(11) NOT NULL DEFAULT '0' COMMENT '终端:授权总数',
			agent_used int(11) DEFAULT '0' COMMENT '终端:已使用目',
			PRIMARY KEY (order_no),
			KEY org_id (org_id)
        ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC COMMENT = '企业:订单授权';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *OrderGrantModel) InsertItem(s *OrderGrantModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateById 通过id修改
func (m *OrderGrantModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Model(m).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *OrderGrantModel) DeleteByConditions(condition map[string]interface{}) bool {
	result := m.Model(m).Where(condition).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *OrderGrantModel) GetItemById(id int64) (item *OrderGrantModel, err error) {
	result := m.Model(m).Where("id=?", id).First(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *OrderGrantModel) GetItemsByIds(ids []int64) (items []*OrderGrantModel, err error) {
	result := m.Model(m).Where("id IN ?", ids).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *OrderGrantModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*OrderGrantModel, err error) {
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

// Count 通过条件统计
func (m *OrderGrantModel) Count(condition map[string]interface{}) int64 {
	var total int64
	result := m.Model(m).Select("id").Where(condition).Count(&total)
	if result.Error == nil {
		return total
	} else {
		return 0
	}
}

// Page 通过条件分页
func (m *OrderGrantModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*OrderGrantModel) {
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
