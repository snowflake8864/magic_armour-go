package order_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 订单表

func CreateOrderFactory() *OrderModel {
	return &OrderModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OrderModel struct {
	model.BaseModel
	OrderNo       string  `json:"order_no" gorm:"primaryKey"`
	OrgId         int64   `json:"org_id"`
	OrgName       string  `json:"org_name"`
	UserName      string  `json:"user_name"`
	UserEmail     string  `json:"user_email"`
	UserPhone     string  `json:"user_phone"`
	ProductId     int64   `json:"product_id"`
	ProductTitle  string  `json:"product_title"`
	GrantType     int     `json:"grant_type"`
	GrantModule   string  `json:"grant_module"`
	SalePrice     float64 `json:"sale_price"`
	GrantDuration int     `json:"grant_duration"`
	BuyNumber     int     `json:"buy_number"`
	UsageNumber   int     `json:"usage_number"`
	TotalAmount   float64 `json:"total_amount"`
	DeclineAmount float64 `json:"decline_amount"`
	RealAmount    float64 `json:"real_amount"`
	OrderStatus   int     `json:"order_status"`
	IsExpire      int     `json:"is_expire"`
	ExpireAt      string  `json:"expire_at"`
	PayMethod     int     `json:"pay_method"`
	PayInfo       string  `json:"pay_info"`
	PayOrderNo    string  `json:"pay_order_no"`
	CreatedAt     string  `json:"created_at"`
	DeclineAt     string  `json:"decline_at"`
	DeclineRemark string  `json:"decline_remark"`
	PayAt         string  `json:"pay_at"`
	CompleteAt    string  `json:"complete_at"`
	CancelAt      string  `json:"cancel_at" gorm:"default:NULL"`
	Source        int     `json:"source"`
}

// 表名
func (m *OrderModel) TableName() string {
	return "o_orders"
}

func (m *OrderModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
          order_no VARCHAR(128) NOT NULL DEFAULT '' COMMENT '订单编号',
          org_id int (11) NOT NULL DEFAULT '0' COMMENT '企业(o_organization.id)',
          org_name varchar (128) NOT NULL DEFAULT '' COMMENT '企业(o_organization.name)',
          user_name varchar(128) NOT NULL DEFAULT '' COMMENT '默认管理用户(o_organization.user_name)',
		  user_email varchar(64) NOT NULL DEFAULT '' COMMENT '默认管理用户(o_organization.user_email)',
		  user_phone varchar(16) NOT NULL DEFAULT '' COMMENT '默认管理用户(o_organization.user_phone)',
          product_id int (11) NOT NULL DEFAULT '0' COMMENT '产品id(p_products.id)',
          product_title varchar (128) NOT NULL DEFAULT '' COMMENT '产品标题(p_products.title)',
          grant_type tinyint (1) NOT NULL DEFAULT '0' COMMENT '授权类型',
          grant_module varchar (1024) NOT NULL DEFAULT '' COMMENT '授权模块(json)',
          sale_price decimal (8, 2) NOT NULL DEFAULT '0.00' COMMENT '产品售价',
          grant_duration int (11) NOT NULL DEFAULT '0' COMMENT '授权时长',
          buy_number int (11) NOT NULL DEFAULT '0' COMMENT '购买数量',
          usage_number  int (11) NOT NULL DEFAULT '0' COMMENT '使用数量',
          total_amount decimal (14, 2) NOT NULL DEFAULT '0.00' COMMENT '总计金额',
          decline_amount decimal (14, 2) NOT NULL DEFAULT '0.00' COMMENT '直降金额',
          real_amount decimal (14, 2) NOT NULL DEFAULT '0.00' COMMENT '真实金额',
          order_status SMALLINT (4) NOT NULL DEFAULT '1' COMMENT '订单状态',
          is_expire TINYINT (1) NOT NULL DEFAULT '0' COMMENT '是否过期',
          expire_at DATE DEFAULT NULL COMMENT '过期时间',
          pay_method tinyint (1) NOT NULL DEFAULT '0' COMMENT '支付方法',
          pay_info text COMMENT '支付信息',
          pay_order_no varchar (128) NOT NULL DEFAULT '' COMMENT '支付订单号(三方)',
          created_at datetime DEFAULT NULL COMMENT '创建时间',
          decline_at datetime DEFAULT NULL COMMENT '直降时间',
          decline_remark varchar (255) NOT NULL DEFAULT '' COMMENT '直降备注',
          pay_at datetime DEFAULT NULL COMMENT '支付时间',
          complete_at datetime DEFAULT NULL COMMENT '完成时间',
          cancel_at datetime DEFAULT NULL COMMENT '取消',
          source tinyint (1) NOT NULL DEFAULT '0' COMMENT '来源',
          PRIMARY KEY (order_no),
          KEY org_expire (org_id, order_status, is_expire),
          KEY org_name (org_name),
          KEY product_id (product_id),
          KEY user_name (user_name),
          KEY created_at (created_at),
          KEY pay_at (pay_at),
          KEY complete_at (complete_at)
        ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = DYNAMIC COMMENT = '企业:订单';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *OrderModel) InsertItem(s *OrderModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteById 通过id删除
func (m *OrderModel) DeleteById(id string) bool {
	result := m.Where("order_no=?", id).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateById 通过id修改
func (m *OrderModel) UpdateById(id string, updates map[string]interface{}) bool {
	result := m.Model(m).Where("order_no=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *OrderModel) GetItemById(id string) (item *OrderModel, err error) {
	result := m.Model(m).Where("order_no=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *OrderModel) GetItemsByIds(ids []string) (items []*OrderModel, err error) {
	result := m.Model(m).Where("order_no IN ?", ids).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *OrderModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*OrderModel, err error) {
	tmp := m.Model(m)
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
	result := tmp.Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// Count 通过条件统计
func (m *OrderModel) Count(condition map[string]interface{}) int64 {
	var total int64
	tmp := m.Model(m).Select("order_no")
	if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Count(&total)
	if result.Error == nil {
		return total
	} else {
		return 0
	}
}

// Page 通过条件分页
func (m *OrderModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*OrderModel) {
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

// CalcByConditions 通过条件获取多条记录
func (m *OrderModel) CalcByConditions(condition map[string]interface{}, column string) (item *OrderModel, err error) {
	tmp := m.Model(m).Select(column)
	if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return nil, result.Error
	}
}
