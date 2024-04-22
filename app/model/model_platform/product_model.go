package model_platform

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
	"time"
)

// 产品表

func CreateProductFactory() *ProductModel {
	return &ProductModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type ProductModel struct {
	model.BaseModel
	Id          int64   `json:"id" gorm:"primaryKey"`
	Title       string  `json:"title"`
	GrantType   int     `json:"grant_type"`
	GrantModule string  `json:"grant_module"`
	Price       float64 `json:"price"`
	Remark      string  `json:"remark"`
	IsPublish   int     `json:"is_publish"`
	CreateId    int64   `json:"create_id"`
	CreateName  string  `json:"create_name"`
	CreatedAt   string  `json:"created_at"`
}

// 表名

func (m *ProductModel) TableName() string {
	return "p_products"
}

func (m *ProductModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			title varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
			grant_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '授权类型',
			grant_module varchar(1024) NOT NULL DEFAULT '' COMMENT '授权模块(json)',
			price decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '销售价格',
			remark varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
			is_publish tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否发布',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			KEY title (title)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='平台:产品';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *ProductModel) InitData() bool {
	curTime := time.Now().Format(variable.DateFormat)

	list := []ProductModel{
		{
			Title:       "幻甲测试授权",
			GrantType:   enum_product.TypeTest,
			GrantModule: "{\"process\":1,\"tamper\":1,\"blackmail\":1,\"port\":1,\"login\":1,\"plug\":1,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       0.00,
			Remark:      "",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲企业版（终端）",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":0,\"blackmail\":0,\"port\":0,\"login\":0,\"plug\":0,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       40.00,
			Remark:      "支持功能：可信防护、外设管理、通讯日志、外联管控、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲企业版（服务器）",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":0,\"blackmail\":0,\"port\":0,\"login\":0,\"plug\":0,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       120.00,
			Remark:      "支持功能：可信防护、外设管理、通讯日志、外联管控、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（终端）",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":0,\"blackmail\":0,\"port\":1\"login\":1,\"plug\":0,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       80.00,
			Remark:      "支持功能：威胁情报模块、可信防护、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（服务器）",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":0,\"blackmail\":0,\"port\":1\"login\":1,\"plug\":0,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       240.00,
			Remark:      "支持功能：威胁情报模块、可信防护、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（终端）+防篡改",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":1,\"blackmail\":0,\"port\":1\"login\":1,\"plug\":0,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       880.00,
			Remark:      "支持功能：威胁情报模块、可信防护、防篡改、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（服务器）+防篡改",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":1,\"blackmail\":0,\"port\":1\"login\":1,\"plug\":0,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       1040.00,
			Remark:      "支持功能：威胁情报模块、可信防护、防篡改、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（终端）+勒索防护",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":0,\"blackmail\":1,\"port\":1\"login\":1,\"plug\":0,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       880.00,
			Remark:      "支持功能：威胁情报模块、可信防护、勒索防护、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（服务器）+勒索防护",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":0,\"blackmail\":1,\"port\":1\"login\":1,\"plug\":0,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       1040.00,
			Remark:      "支持功能：威胁情报模块、可信防护、勒索防护、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（终端）+动态封堵",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":0,\"blackmail\":0,\"port\":1\"login\":1,\"plug\":1,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       880.00,
			Remark:      "支持功能：威胁情报模块、可信防护、动态封堵、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（服务器）+动态封堵",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":0,\"blackmail\":0,\"port\":1\"login\":1,\"plug\":1,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       1040.00,
			Remark:      "支持功能：威胁情报模块、可信防护、动态封堵、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（终端）+防篡改+勒索防护",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":1,\"blackmail\":1,\"port\":1\"login\":1,\"plug\":0,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       1680.00,
			Remark:      "支持功能：威胁情报模块、可信防护、防篡改、勒索防护、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（服务器）+防篡改+勒索防护",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":1,\"blackmail\":1,\"port\":1\"login\":1,\"plug\":0,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       1840.00,
			Remark:      "支持功能：威胁情报模块、可信防护、防篡改、勒索防护、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（终端）+防篡改+勒索防护+动态封堵",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":1,\"blackmail\":1,\"port\":1\"login\":1,\"plug\":1,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       2480.00,
			Remark:      "支持功能：威胁情报模块、可信防护、防篡改、勒索防护、动态封堵、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
		{
			Title:       "幻甲旗舰版（服务器）+防篡改+勒索防护+动态封堵",
			GrantType:   enum_product.TypeGold,
			GrantModule: "{\"process\":1,\"tamper\":1,\"blackmail\":1,\"port\":1\"login\":1,\"plug\":1,\"usb\":1,\"outreach\":1,\"syslog\":1}",
			Price:       2640.00,
			Remark:      "支持功能：威胁情报模块、可信防护、防篡改、勒索防护、动态封堵、端口虚开、口令破解防护、外设管理、外联管控、通讯日志、安全报告、安全通知等。",
			IsPublish:   enums.IsYes,
			CreatedAt:   curTime,
		},
	}
	result := m.Create(&list)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *ProductModel) InsertItem(s *ProductModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateById 通过id修改
func (m *ProductModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Model(m).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteById 通过id删除
func (m *ProductModel) DeleteById(id int64) bool {
	result := m.Model(m).Where("id=?", id).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *ProductModel) GetItemById(id int64) (item *ProductModel, err error) {
	result := m.Model(m).Where("id=?", id).First(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *ProductModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*ProductModel, err error) {
	whereSQL, vals, _ := m.WhereBuild(condition)
	tmp := m.Model(m).Where(whereSQL, vals)
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
func (m *ProductModel) Count(condition map[string]interface{}) int64 {
	var total int64
	whereSQL, vals, _ := m.WhereBuild(condition)
	result := m.Model(m).Select("id").Where(whereSQL, vals).Count(&total)
	if result.Error == nil {
		return total
	} else {
		return 0
	}
}

// Page 通过条件分页
func (m *ProductModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*ProductModel) {
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
