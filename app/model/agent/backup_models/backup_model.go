package backup_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 备份表

func CreateBackupFactory() *BackupModel {
	return &BackupModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type BackupModel struct {
	baseBackup
	model.BaseModel
	Id               string `json:"id" gorm:"primaryKey"`
	OrgId            int64  `json:"org_id"`
	AgentId          string `json:"agent_id"`
	Name             string `json:"name"`
	Remark           string `json:"remark"`
	SourcePath       string `json:"source_path"`
	StorageDir       string `json:"storage_dir"`
	FileName         string `json:"file_name"`
	FileSize         string `json:"file_size"`
	BackupStatus     int    `json:"backup_status"`
	BackupId         int64  `json:"backup_id"`
	BackupName       string `json:"backup_name"`
	BackupStartAt    string `json:"backup_start_at" gorm:"default:NULL"`
	BackupEndAt      string `json:"backup_end_at" gorm:"default:NULL"`
	BackupFailReason string `json:"backup_fail_reason"`
	RollbackStatus   int    `json:"rollback_status"`
	RollbackId       int64  `json:"rollback_id"`
	RollbackName     string `json:"rollback_name"`
	RollbackStartAt  string `json:"rollback_start_at" gorm:"default:NULL"`
	RollbackEndAt    string `json:"rollback_end_at"  gorm:"default:NULL"`
	RollbackReason   string `json:"rollback_fail_reason"`
	CreatedAt        string `json:"created_at"`
}

// 表名

func (m *BackupModel) TableName() string {
	return "a_backups"
}

func (m *BackupModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT '设备:hash',
			name varchar(128)  NOT NULL DEFAULT '' COMMENT '名称',
			remark varchar(128)  NOT NULL DEFAULT '' COMMENT '备注',
			source_path varchar(512)  NOT NULL DEFAULT '' COMMENT '原路径',
			storage_dir varchar(512)  NOT NULL DEFAULT '' COMMENT '存储路径',
			file_name varchar(128)  NOT NULL DEFAULT '' COMMENT '备份文件名',
			file_size varchar(32)  NOT NULL DEFAULT '' COMMENT '备份文件大小',
			backup_status tinyint NOT NULL DEFAULT '1' COMMENT '备份状态(1默认 2进行中 3成功 4失败)',
			backup_id int NOT NULL DEFAULT '0' COMMENT '备份用户(user.id)',
			backup_name varchar(128)  NOT NULL DEFAULT '' COMMENT '备份用户(user.username)',
			backup_start_at datetime NULL COMMENT '备份开始时间',
			backup_end_at datetime NULL COMMENT '备份结束时间',
		    backup_fail_reason varchar(512)  NOT NULL DEFAULT '' COMMENT '备份失败原因',
			rollback_status tinyint NOT NULL DEFAULT '1' COMMENT '备份状态(1默认 2进行中 3成功 4失败)',
			rollback_id int NOT NULL DEFAULT '0' COMMENT '回滚用户(user.id)',
			rollback_name varchar(128)  NOT NULL DEFAULT '' COMMENT '回滚用户(user.username)',
			rollback_start_at datetime NULL COMMENT '回滚开始时间',
			rollback_end_at datetime NULL COMMENT '回滚结束时间',
			rollback_fail_reason varchar(512)  NOT NULL DEFAULT '' COMMENT '回滚失败原因',
			created_at datetime NOT NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id,agent_id, id) 
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

func (m *BackupModel) InsertItem(s *BackupModel) bool {

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

// InsertBatch 插入多条
func (m *BackupModel) InsertBatch(list *[]BackupModel) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(list)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return false
	}
}

// UpdateByConditions 通过条件修改
func (m *BackupModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) (bool, int64) {
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
func (m *BackupModel) DeleteByConditions(condition SearchConditions) bool {
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

// GetItemByOrgAndAgentAndHash 通过企业id\终端id\hash获取一条记录
func (m *BackupModel) GetItemByOrgAndAgentAndHash(orgId int64, agentId, hash string) (item *BackupModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND agent_id=? AND id=?", orgId, agentId, hash).Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item
	} else {
		return nil
	}
}

// GetItemByConditions 通过id获取一条记录
func (m *BackupModel) GetItemByConditions(condition SearchConditions) (item *BackupModel) {
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
func (m *BackupModel) GetListByConditions(condition SearchConditions) (items []*BackupModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *BackupModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*BackupModel) {
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
func (m *BackupModel) Count(condition SearchConditions) int64 {
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
func (m *BackupModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*BackupModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}

// CalcByConditionsGroupAgent 通过条件获取多条记录且按终端分组
func (m *BackupModel) CalcByConditionsGroupAgent(condition SearchConditions, column string) (items []*Calc) {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Select(column)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Group("agent_id").Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items
	} else {
		return nil
	}
}

// CalcByConditionsGroupAgentReturnMap 通过条件获取多条记录且按终端分组
func (m *BackupModel) CalcByConditionsGroupAgentReturnMap(condition SearchConditions, column string) map[string]*Calc {
	var result = make(map[string]*Calc)

	items := m.CalcByConditionsGroupAgent(condition, column)
	if items == nil || len(items) == 0 {
		return result
	}
	for _, item := range items {
		result[item.AgentId] = item
	}

	return result
}
