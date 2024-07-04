package organization_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 组织架构表

func CreateOrganizationFactory() *OrganizationModel {
	return &OrganizationModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OrganizationModel struct {
	baseOrg
	model.BaseModel
	Id                int64  `json:"id" gorm:"primaryKey"`
	Code              string `json:"code"`
	Name              string `json:"name"`
	Type              int    `json:"type"`
	UserName          string `json:"user_name"`
	UserEmail         string `json:"user_email"`
	UserPhone         string `json:"user_phone"`
	Address           string `json:"address"`
	Icon              string `json:"icon"`
	Status            int    `json:"status"`
	IsClean           int    `json:"is_clean"`
	IsDel             int    `json:"is_del"`
	Remark            string `json:"remark"`
	LoginPhoneCheck   int    `json:"login_phone_check"`
	HuanshiSwitch     int    `json:"huanshi_switch"`
	TrusteeshipSwitch int    `json:"trusteeship_switch"`
	Source            int    `json:"source"`
	CreateId          int64  `json:"create_id"`
	CreateName        string `json:"create_name"`
	CreatedAt         string `json:"created_at"`
}

// 表名
func (m *OrganizationModel) TableName() string {
	return "o_organizations"
}

func (m *OrganizationModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
		 id int(11) unsigned NOT NULL AUTO_INCREMENT,
		 code varchar(64) DEFAULT NULL COMMENT '标识码',
		 name varchar(128) DEFAULT NULL COMMENT '名称',
		 type tinyint(1) NOT NULL DEFAULT '1' COMMENT '类型:(1企业 2个人)',
		 user_name varchar(128) NOT NULL DEFAULT '' COMMENT '默认管理用户(o_user.name)',
		 user_email varchar(64) NOT NULL DEFAULT '' COMMENT '默认管理用户(o_user.email)',
		 user_phone varchar(16) NOT NULL DEFAULT '' COMMENT '默认管理用户(o_user.phone)',
		 address varchar(512) NOT NULL DEFAULT '' COMMENT '地址',
		 icon varchar(512) NOT NULL DEFAULT '' COMMENT '头像',
		 status tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态:(0审核中,1正常,2锁定)',
		 is_clean TINYINT(1) DEFAULT 0 NOT NULL COMMENT '是否清理数据',
		 is_del tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
		 remark varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
		 login_phone_check tinyint(1) NOT NULL DEFAULT '0' COMMENT '登录手机验证(0否,1是)',
		 huanshi_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '幻势开关',
		 trusteeship_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '托管开关',
		 source tinyint(1) NOT NULL DEFAULT '0' COMMENT '来源(1自注册 2平台添加)',
		 create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(p_admin.id)',
		 create_name varchar(128) NOT NULL DEFAULT '' COMMENT '创建人(p_admin.name)',
		 created_at datetime DEFAULT NULL,
		 updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		  PRIMARY KEY (id),
		  UNIQUE KEY code (code),
		  KEY name (name),
		  KEY user_name (user_name),
		  KEY user_email (user_email),
		  KEY user_phone (user_phone)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='企业:信息表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *OrganizationModel) InsertItem(s *OrganizationModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateById 通过id修改
func (m *OrganizationModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Model(m).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteById 通过id删除
func (m *OrganizationModel) DeleteById(id int64) bool {
	result := m.Where("id=?", id).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *OrganizationModel) DeleteByConditions(condition SearchConditions) bool {
	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *OrganizationModel) GetItemById(id int64) (item *OrganizationModel, err error) {
	result := m.Model(m).Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

func (m *OrganizationModel) GetItemByName(name string) (item *OrganizationModel, err error) {
	result := m.Model(m).Where("name=?", name).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *OrganizationModel) GetItemsByIds(ids []int64) (items []*OrganizationModel, err error) {
	result := m.Model(m).Where("id IN ?", ids).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// GetListByConditions 通过条件获取多条记录
func (m *OrganizationModel) GetListByConditions(condition SearchConditions) (items []*OrganizationModel) {
	items, _ = m.GetItemsByConditions(condition, 0, 0, "")
	return items
}

// GetItemsByConditions 通过条件获取多条记录
func (m *OrganizationModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*OrganizationModel, err error) {
	tmp := m.Model(m)
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
	result := tmp.Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// Count 通过条件统计
func (m *OrganizationModel) Count(condition SearchConditions) int64 {
	var total int64
	tmp := m.Model(m).Select("id")
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
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
func (m *OrganizationModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*OrganizationModel) {
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
