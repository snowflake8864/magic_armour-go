package proc_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

// 模块运行中表

func CreateDllRunningFactory() *DllRunningModel {
	return &DllRunningModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type DllRunningModel struct {
	baseProc
	model.BaseModel
	Id        string `json:"id" gorm:"primaryKey"`
	OrgId     int64  `json:"org_id"`
	AgentId   string `json:"agent_id"`
	DllHash   string `json:"dll_hash"`
	ProcId    int    `json:"proc_id"`
	DllDir    string `json:"dll_dir"`
	ProcHash  string `json:"proc_hash"`
	ProcDir   string `json:"proc_dir"`
	CreatedAt string `json:"created_at"`
}

// 表名

func (m *DllRunningModel) TableName() string {
	return "a_dll_running"
}

func (m *DllRunningModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			dll_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '模块:hash',
			proc_id int NOT NULL DEFAULT '0'  COMMENT '进程:id',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			dll_dir varchar(512)  NOT NULL DEFAULT '' COMMENT '模块:路径',
            proc_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '进程:hash',
            proc_dir varchar(512)  NOT NULL DEFAULT '' COMMENT '进程:路径',
			created_at datetime NULL,
			INDEX index_dll_hash (dll_hash) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id, dll_hash, proc_id) 
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

func (m *DllRunningModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *DllRunningModel) InsertItem(s *DllRunningModel) bool {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(s)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByOrgAndAgent 通过企业id\终端id删除
func (m *DllRunningModel) DeleteByOrgAndAgent(orgId int64, agentId string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		tmp = tmp.Where("org_id=? AND agent_id=? ", orgId, agentId)
		return tmp.Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *DllRunningModel) DeleteByConditions(condition SearchConditions) bool {
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
func (m *DllRunningModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*DllAndRunning) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {

		t1As := m.TableName() + "."
		t2Name := CreateDllFactory().TableName()
		t2As := "t2."

		columns := []string{
			t1As + "org_id",
			t1As + "agent_id",
			t1As + "dll_hash",
			t1As + "id",
			t1As + "dll_dir",
			t1As + "created_at",
			t2As + "status",
			t2As + "sync_status_at",
			t2As + "self_status",
			t2As + "threaten_status",
			t2As + "auto_white_switch",
			t2As + "introduce",
			t2As + "copyright",
			t2As + "signature",
			t2As + "attribute",
		}
		tmp := tx.Model(m).
			Joins(m.getJoinStr(t1As, t2Name, t2As))

		tmp = tmp.Select(strings.Join(columns, ","))
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2WhereForRunning(condition, t1As, t2As)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		if limit > 0 {
			tmp = tmp.Offset(offset).Limit(limit)
		}
		if !string_helper.IsEmpty(order) {
			tmp = tmp.Order(t1As + order)
		} else {
			tmp = tmp.Order(t1As + "created_at ASC")
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
func (m *DllRunningModel) Count(condition SearchConditions) int64 {
	var total int64
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		t1As := m.TableName() + "."
		t2Name := CreateDllFactory().TableName()
		t2As := "t2."

		tmp := tx.Model(m).
			Joins(m.getJoinStr(t1As, t2Name, t2As))
		tmp = tmp.Select(t1As + "id")
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2WhereForRunning(condition, t1As, t2As)); !string_helper.IsEmpty(whereSQL) {
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
func (m *DllRunningModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*DllAndRunning) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
