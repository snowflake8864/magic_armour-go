package software_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 软件表

func CreateSoftwareFactory() *SoftwareModel {
	return &SoftwareModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type SoftwareModel struct {
	baseSoftware
	model.BaseModel
	Id           int64  `json:"id" gorm:"primaryKey"`
	UploadId     int64  `json:"upload_id"`
	Ver          string `json:"ver"`
	PlatformType int    `json:"platform_type"`
	OsType       int    `json:"os_type"`
	OsVer        string `json:"os_ver"`
	IsPublish    int    `json:"is_publish"`
	Remark       string `json:"remark"`
	StoragePath  string `json:"storage_path"`
	FullPath     string `json:"full_path"`
	CreateId     int64  `json:"create_id"`
	CreateName   string `json:"create_name"`
	CreatedAt    string `json:"created_at"`
}

// 表名

func (m *SoftwareModel) TableName() string {
	return "p_software"
}

func (m *SoftwareModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			upload_id int(11) NOT NULL DEFAULT '0' COMMENT '(upload.id)',
            ver varchar(64)  NOT NULL DEFAULT '' COMMENT '版本号(0.0.0.0)',
			platform_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '平台类型(win|linxu, 统一定义)',
			os_type smallint(4) NOT NULL DEFAULT '0' COMMENT '系统类别(参考ubuntu|centos等,统一定义)',
			os_ver varchar(512)  NOT NULL DEFAULT '' COMMENT '系统版本(统一定义)',
			is_publish tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否发布',
            remark text  NOT NULL COMMENT '备注',
			storage_path varchar(512)  NOT NULL DEFAULT '' COMMENT '存储相对路径(项目下)',
			full_path varchar(1024)  NOT NULL DEFAULT '' COMMENT '完整路径',
            create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			UNIQUE KEY(id),
			KEY index_ver (ver),
			KEY index_platform_type (platform_type),
			KEY index_os_type (os_type)
		)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='软件表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *SoftwareModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *SoftwareModel) InsertItem(s *SoftwareModel) bool {

	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *SoftwareModel) GetItemById(id int64) (item *SoftwareModel, err error) {
	result := m.Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return nil, result.Error
	}
}

func (m *SoftwareModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Table(m.TableName()).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *SoftwareModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) bool {
	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteById 通过id删除
func (m *SoftwareModel) DeleteById(id int64) bool {
	result := m.Model(m).Where("id=?", id).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *SoftwareModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*SoftwareModel) {
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
func (m *SoftwareModel) Count(condition SearchConditions) int64 {
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
func (m *SoftwareModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*SoftwareModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
