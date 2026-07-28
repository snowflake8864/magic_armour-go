package antivirus_engine_task_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 上线日志表

func CreateUpgradeFailFactory() *UpgradeFailModel {
	return &UpgradeFailModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type UpgradeFailModel struct {
	baseOnline
	model.BaseModel
	Id          int64  `json:"id" gorm:"primaryKey"`
	OrgId       int64  `json:"org_id"`
	AgentId     string `json:"agent_id"`
	HostName    string `json:"host_name"`
	MacId       string `json:"mac_id"`
	AgentIp     string `json:"agent_ip"`
	NetIp       string `json:"net_ip"`
	OsName      string `json:"os_name"`
	IsOnline    int    `json:"is_online"`
	DeptId      int64  `json:"dept_id"`
	Ver         string `json:"ver"`
	Alias       string `json:"alias"`
	ContactName string `json:"contact_name"`
	ContactTel  string `json:"contact_tel"`
	Remark      string `json:"remark"`
	CreatedAt   string `json:"created_at"`
}

// 表名
func (m *UpgradeFailModel) TableName() string {
	return "a_antivirus_engine_upgrade_failed"
}

func (m *UpgradeFailModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			task_id int(11) unsigned NOT NULL  COMMENT '任务ID',
			exec_batch int(11) NOT NULL DEFAULT '0' COMMENT '执行批次',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			created_at datetime NULL,
			id BIGINT NOT NULL DEFAULT '0',
			dept_id int NOT NULL DEFAULT '0' COMMENT '组织架构编号(department.id)',
			reason tinyint(1) NOT NULL DEFAULT '0' COMMENT '原因',
			remark varchar(128)  NOT NULL DEFAULT '' COMMENT '备注'
		)
		DUPLICATE  KEY(org_id, task_id, exec_batch, agent_id, created_at)
		DISTRIBUTED BY HASH(org_id) BUCKETS AUTO 
		PROPERTIES (
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

func (m *UpgradeFailModel) AlterTable(ver string) bool {
	sql := `ALTER TABLE ` + m.TableName()
	switch ver {
	default:
		return true
	}
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *UpgradeFailModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *UpgradeFailModel) InsertItem(s *UpgradeFailModel) bool {

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
func (m *UpgradeFailModel) InsertBatch(list *[]UpgradeFailModel) bool {
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

// DeleteByOrg 通过企业i删除
func (m *UpgradeFailModel) DeleteByOrg(orgId int64) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *UpgradeFailModel) DeleteByConditions(condition SearchConditions) bool {
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

// GetItemsByConditions 通过条件获取多条记录
func (m *UpgradeFailModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*UpgradeFailModel) {
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
			tmp = tmp.Order("created_at DESC")
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
func (m *UpgradeFailModel) Count(condition SearchConditions) int64 {
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
func (m *UpgradeFailModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*UpgradeFailModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
