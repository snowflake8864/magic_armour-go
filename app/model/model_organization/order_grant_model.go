package model_organization

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
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
func (m *OrderGrantModel) UpdateById(id string, updates map[string]interface{}) bool {
	result := m.Model(m).Where("order_no=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *OrderGrantModel) GetItemById(id string) (item *OrderGrantModel, err error) {
	result := m.Model(m).Where("order_no=?", id).First(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过ids获取多条记录
func (m *OrderGrantModel) GetItemsByIds(ids []string) (items []*OrderGrantModel, err error) {
	result := m.Model(m).Where("order_no IN ?", ids).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}
