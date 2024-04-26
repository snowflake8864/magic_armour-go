package model_platform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 软件表

func CreateSoftwareFactory() *SoftwareModel {
	return &SoftwareModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type SoftwareModel struct {
	model.BaseModel
	Id           int64  `json:"id" gorm:"primaryKey"`
	UploadId     int64  `json:"upload_id"`
	Ver          string `json:"ver"`
	PlatformType int    `json:"platform_type"`
	OsType       int    `json:"os_type"`
	OsVer        string `json:"os_ver"`
	IsPublish    int    `json:"is_publish"`
	Remark       string `json:"remark"`
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
            create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			UNIQUE KEY(id),
			KEY index_ver (ver),
			KEY index_platform_type (platform_type),
			KEY index_os_type (os_type)
		)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='软甲表';
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

func (m *SoftwareModel) GetItemById(id int64) (*SoftwareModel, error) {
	result := m.Where("id=?", id).First(m)
	if result.Error == nil {
		return m, nil
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
func (m *SoftwareModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*SoftwareModel, err error) {
	tmp := m.Model(m).Where(condition)
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
func (m *SoftwareModel) Count(condition map[string]interface{}) int64 {
	var total int64
	result := m.Model(m).Select("id").Where(condition).Count(&total)
	if result.Error == nil {
		return total
	} else {
		return 0
	}
}

// Page 通过条件分页
func (m *SoftwareModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*SoftwareModel) {
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
