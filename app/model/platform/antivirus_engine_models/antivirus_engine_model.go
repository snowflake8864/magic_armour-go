package antivirus_engine_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 杀毒引擎表

func CreateEngineFactory() *EngineModel {
	return &EngineModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type EngineModel struct {
	baseEngine
	model.BaseModel
	Id           int64  `json:"id" gorm:"primaryKey"`
	UploadId     int64  `json:"upload_id"`
	PlatformType int    `json:"platform_type"`
	ArchType     int    `json:"arch_type"`
	EngineType   int    `json:"engine_type"`
	Ver          string `json:"ver"`
	IsPublish    int    `json:"is_publish"`
	IsNewest     int    `json:"is_newest"`
	Remark       string `json:"remark"`
	FileHash     string `json:"file_hash"`
	FileSize     int64  `json:"file_size"`
	StoragePath  string `json:"storage_path"`
	FullPath     string `json:"full_path"`
	CreateId     int64  `json:"create_id"`
	CreateName   string `json:"create_name"`
	CreatedAt    string `json:"created_at"`
}

// 表名

func (m *EngineModel) TableName() string {
	return "p_antivirus_engines"
}

func (m *EngineModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            id int(11) unsigned NOT NULL AUTO_INCREMENT,
			upload_id int(11) NOT NULL DEFAULT '0' COMMENT '(upload.id)',
			platform_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '平台类型(win|linux, 统一定义)',
			arch_type tinyint(1) NOT NULL DEFAULT '1' COMMENT 'CPU架构(参考amd64|arm64等,统一定义)',
			engine_type smallint(4)  unsigned NOT NULL DEFAULT '0' COMMENT '引擎类型',
            ver varchar(64)  NOT NULL DEFAULT '' COMMENT '版本号',
			is_publish tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否发布',
			is_newest tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否最新',
            remark text  NOT NULL COMMENT '备注',
			file_size int unsigned NOT NULL DEFAULT '0'  COMMENT '文件大小',
  			file_hash varchar(32) NOT NULL DEFAULT ''  COMMENT '文件哈希值',
			storage_path varchar(512)  NOT NULL DEFAULT '' COMMENT '存储相对路径(项目下)',
			full_path varchar(1024)  NOT NULL DEFAULT '' COMMENT '完整路径',
            create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(admin.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(admin.username)',
			created_at datetime NOT NULL,
			UNIQUE KEY(id),
			KEY index_ver (ver),
			KEY index_platform_type (platform_type)
		)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='杀毒引擎表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *EngineModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *EngineModel) InsertItem(s *EngineModel) bool {

	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *EngineModel) GetItemById(id int64) (item *EngineModel, err error) {
	result := m.Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return nil, result.Error
	}
}

func (m *EngineModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Table(m.TableName()).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *EngineModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) bool {
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
func (m *EngineModel) DeleteById(id int64) bool {
	result := m.Model(m).Where("id=?", id).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemByConditions 通过条件获一条记录
func (m *EngineModel) GetItemByConditions(condition SearchConditions) (item *EngineModel) {
	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}

	result := tmp.Take(&item)
	if result.Error == nil {
		return item
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return nil
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *EngineModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*EngineModel) {
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
func (m *EngineModel) Count(condition SearchConditions) int64 {
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
func (m *EngineModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*EngineModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
