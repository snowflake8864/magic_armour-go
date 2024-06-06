package tamper_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

//防篡改-全局

func CreateTamperGlobalFactory() *TamperGlobalModel {
	return &TamperGlobalModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type TamperGlobalModel struct {
	baseTamper
	model.BaseModel
	Id            int64  `json:"id" gorm:"primaryKey"`
	OrgId         int64  `json:"org_id"`
	PlatformType  int    `json:"platform_type"`
	DeptId        int64  `json:"dept_id"`
	Type          int    `json:"type"`
	Dir           string `json:"dir"`
	ProtectRw     int    `json:"protect_rw"`
	IncludeSuffix string `json:"include_suffix"`
	ExcludeSuffix string `json:"exclude_suffix"`
	ExcludeFolder string `json:"exclude_folder"`
	IsExtend      int    `json:"is_extend"`
	Processes     string `json:"processes"`
	OpenSwitch    int    `json:"open_switch"`
	CreateId      int64  `json:"create_id"`
	CreateName    string `json:"create_name"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// 表名

func (m *TamperGlobalModel) TableName() string {
	return "a_tamper_global"
}

func (m *TamperGlobalModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			platform_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '平台类型(统一定义)',
			id BIGINT NOT NULL DEFAULT '0',
			dept_id int NOT NULL DEFAULT '0' COMMENT '组织架构编号(department.id)',
			type tinyint NOT NULL DEFAULT '0' COMMENT '类型(1文件夹 2单文件)',
			dir varchar(512)  NOT NULL DEFAULT '' COMMENT '文件路径',
			protect_rw int NOT NULL DEFAULT '0' COMMENT '保护权限(1读,2写,4删除,8重命名,16新建)',
			include_suffix varchar(256)  NOT NULL DEFAULT '' COMMENT '包含文件类型(|隔开)',
			exclude_suffix varchar(256)  NOT NULL DEFAULT '' COMMENT '排除文件类型(|隔开)',
			exclude_folder varchar(1024)  NOT NULL DEFAULT '' COMMENT '排除文件夹(|隔开)',
			is_extend tinyint NOT NULL DEFAULT '0' COMMENT '是否向下继承',
			processes string  NOT NULL DEFAULT ''  COMMENT '信任进程json',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			updated_at datetime NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, platform_type, id) 
		DISTRIBUTED BY HASH(org_id) BUCKETS AUTO
		PROPERTIES (
			'enable_unique_key_merge_on_write' = 'true',
			'replication_allocation' = 'tag.location.default: 1'
		);
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *TamperGlobalModel) InsertItem(s *TamperGlobalModel) bool {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(s)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// UpdateByConditions 通过条件修改
func (m *TamperGlobalModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) (bool, int64) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// DeleteByConditions 通过条件删除
func (m *TamperGlobalModel) DeleteByConditions(condition SearchConditions) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemByConditions 通过id获取一条记录
func (m *TamperGlobalModel) GetItemByConditions(condition SearchConditions) (item *TamperGlobalModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return nil
	}
}

// GetListByConditions 获取符合条件的全部记录
func (m *TamperGlobalModel) GetListByConditions(condition SearchConditions) (items []*TamperGlobalModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *TamperGlobalModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*TamperGlobalModel) {
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

// Count 通过条件统计
func (m *TamperGlobalModel) Count(condition SearchConditions) int64 {
	var total int64
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Select("id")
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Count(&total)
	})
	result := m.Raw(m.TrimSql(sql)).Count(&total)
	if result.Error == nil {
		return total
	} else {
		return 0
	}
}

// Page 通过条件分页
func (m *TamperGlobalModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*TamperGlobalModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
