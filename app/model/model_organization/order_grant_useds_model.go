package model_organization

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 订单授权使用表

func CreateOrderGrantUsedFactory() *OrderGrantUsedModel {
	return &OrderGrantUsedModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OrderGrantUsedModel struct {
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	OrderNo    string `json:"order_no"`
	UsedNumber int    `json:"used_number"`
	UsedAgents string `json:"used_agents"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *OrderGrantUsedModel) TableName() string {
	return "o_order_grant_used"
}

func (m *OrderGrantUsedModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int (11) unsigned NOT NULL AUTO_INCREMENT,
			org_id int (11) NOT NULL DEFAULT '0' COMMENT '企业(o_organization.id)',
			order_no varchar (128) NOT NULL DEFAULT '' COMMENT '订单编号',
			used_number int (11) NOT NULL DEFAULT '0' COMMENT '使用数目',
			used_agents text COMMENT '使用终端集合',
			create_id int (11) NOT NULL DEFAULT '0' COMMENT '创建人(o_user.id)',
			create_name varchar (64) NOT NULL DEFAULT '' COMMENT '创建人(o_user.username)',
			create_at datetime DEFAULT NULL COMMENT '创建时间',
			PRIMARY KEY (id),
			KEY org_order (org_id, order_no)
		) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC COMMENT = '企业:订单授权使用记录';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *OrderGrantUsedModel) InsertItem(s *OrderGrantUsedModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *OrderGrantUsedModel) DeleteByConditions(condition map[string]interface{}) bool {
	result := m.Model(m).Where(condition).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *OrderGrantUsedModel) GetItemById(id int64) (item *OrderGrantUsedModel, err error) {
	result := m.Model(m).Where("id=?", id).First(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *OrderGrantUsedModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*OrderGrantUsedModel, err error) {
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
