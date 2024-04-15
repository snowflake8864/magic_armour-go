package model_organization

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 订单日志表

func CreateOrderLogFactory() *OrderLogModel {
	return &OrderLogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OrderLogModel struct {
	model.BaseModel
	Id            int64  `json:"id" gorm:"primaryKey"`
	OrgId         int64  `json:"org_id"`
	OrderNo       string `json:"order_no"`
	OrderStatus   int    `json:"order_status"`
	OperateAction int    `json:"operate_action"`
	OperateSource int    `json:"operate_source"`
	OperateId     int64  `json:"operate_id"`
	OperateName   string `json:"operate_name"`
	CreatedAt     string `json:"created_at"`
}

// 表名
func (m *OrderLogModel) TableName() string {
	return "o_order_logs"
}

func (m *OrderLogModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int (11) unsigned NOT NULL AUTO_INCREMENT,
			org_id int (11) NOT NULL DEFAULT '0' COMMENT '企业(o_organization.id)',
			order_no varchar (128) NOT NULL DEFAULT '' COMMENT '订单编号',
			order_status SMALLINT (4) NOT NULL DEFAULT '1' COMMENT '订单状态',
			operate_action SMALLINT (4) NOT NULL DEFAULT '1'  COMMENT '操作:行为',
			operate_source tinyint (1) NOT NULL DEFAULT '0' COMMENT '操作:来源',
			operate_id int (11) NOT NULL DEFAULT '0' COMMENT '操作用户:编号',
			operate_name varchar (128) NOT NULL DEFAULT '' COMMENT '操作用户:名称',
			created_at datetime DEFAULT NULL,
			PRIMARY KEY (id),
			KEY org_order (org_id,order_no)
		) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC COMMENT = '企业:订单日志';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *OrderLogModel) InsertItem(s *OrderLogModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateById 通过id修改
func (m *OrderLogModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Model(m).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *OrderLogModel) DeleteByConditions(condition map[string]interface{}) bool {
	result := m.Model(m).Where(condition).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *OrderLogModel) GetItemById(id int64) (item *OrderLogModel, err error) {
	result := m.Model(m).Where("id=?", id).First(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *OrderLogModel) GetItemsByIds(ids []int64) (items []*OrderLogModel, err error) {
	result := m.Model(m).Where("id IN ?", ids).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *OrderLogModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*OrderLogModel, err error) {
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
func (m *OrderLogModel) Count(condition map[string]interface{}) int64 {
	var total int64
	result := m.Model(m).Select("id").Where(condition).Count(&total)
	if result.Error == nil {
		return total
	} else {
		return 0
	}
}

// Page 通过条件分页
func (m *OrderLogModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*OrderLogModel) {
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
