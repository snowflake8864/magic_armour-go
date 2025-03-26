package role_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 角色表

func CreateRoleFactory() *RoleModel {
	return &RoleModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type RoleModel struct {
	baseRole
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	Name       string `json:"name"`
	Remark     string `json:"remark"`
	Menus      string `json:"menus"`
	IsDefault  int    `json:"is_default"`
	CreateId   int64  `json:"create_id"`
	CreateName string `son:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名
func (m *RoleModel) TableName() string {
	return "o_roles"
}

func (m *RoleModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			name varchar(128) NOT NULL DEFAULT '' COMMENT '名称',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '邮箱',
			menus text DEFAULT NULL COMMENT '邮箱',
			is_default tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认(0否,1是)',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			KEY index_org (org_id),
            KEY index_name (name)
		)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='企业:角色表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *RoleModel) InsertItem(s *RoleModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// InsertBatch 插入多条
func (m *RoleModel) InsertBatch(list *[]RoleModel) bool {
	result := m.Create(list)
	if result.RowsAffected > 0 {
		return true
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *RoleModel) GetItemById(id int64) (item *RoleModel) {
	result := m.Model(m).Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item
	} else {
		return item
	}
}

// GetItemByConditions 通过条件获取一条记录
func (m *RoleModel) GetItemByConditions(condition SearchConditions) (item *RoleModel) {
	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Take(&item)
	if result.Error == nil {
		return item
	} else {
		return item
	}
}

// UpdateByConditions 通过条件修改
func (m *RoleModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) (bool, int64) {
	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Updates(updates)
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// DeleteByConditions 通过条件删除
func (m *RoleModel) DeleteByConditions(condition SearchConditions) bool {
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

// GetListByConditions 通过条件获取多条记录
func (m *RoleModel) GetListByConditions(condition SearchConditions) (items []*RoleModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *RoleModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*RoleModel) {
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
		return items
	} else {
		return nil
	}
}

// Count 通过条件统计
func (m *RoleModel) Count(condition SearchConditions) int64 {
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
func (m *RoleModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*RoleModel) {
	total = m.Count(condition)
	if total == 0 {
		return total, nil
	}
	items = m.GetItemsByConditions(condition, offset, limit, order)
	if len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
