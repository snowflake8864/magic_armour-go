package proc_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

// 进程运行中表

func CreateProcessRunningFactory() *ProcessRunningModel {
	return &ProcessRunningModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ProcessRunningModel struct {
	baseProc
	model.BaseModel
	Id        string `json:"id" gorm:"primaryKey"`
	OrgId     int64  `json:"org_id"`
	AgentId   string `json:"agent_id"`
	ProcHash  string `json:"proc_hash"`
	ProcId    int    `json:"proc_id"`
	ProcDir   string `json:"proc_dir"`
	ProcUser  string `json:"proc_user"`
	DllNumber int    `json:"dll_number"`
	DllStatus int    `json:"dll_status"`
	DllList   string `json:"dll_list"`
	CreatedAt string `json:"created_at"`
}

type ProcessRunningDll struct {
	Hash      string `json:"hash"`
	Dir       string `json:"dir"`
	Introduce string `json:"introduce"`
	Copyright string `json:"copyright"`
	Signature string `json:"signature"`
	Attribute string `json:"attribute"`
}

// 表名

func (m *ProcessRunningModel) TableName() string {
	return "a_process_running"
}

func (m *ProcessRunningModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			proc_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '进程:hash',
			proc_id int NOT NULL DEFAULT '0'  COMMENT '进程:id',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
            proc_dir varchar(512)  NOT NULL DEFAULT '' COMMENT '进程:路径',
            proc_user varchar(128)  NOT NULL DEFAULT '' COMMENT '进程:用户',
            dll_number int  NOT NULL DEFAULT '0' COMMENT '模块:数量',
            dll_status tinyint  NOT NULL DEFAULT '0' COMMENT '模块:状态',
            dll_list string  NOT NULL DEFAULT '' COMMENT '模块:列表',
			created_at datetime NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id, proc_hash, proc_id) 
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

func (m *ProcessRunningModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *ProcessRunningModel) InsertItem(s *ProcessRunningModel) bool {

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
func (m *ProcessRunningModel) DeleteByOrgAndAgent(orgId int64, agentId string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		tmp = tmp.Where("org_id=? AND agent_id=?", orgId, agentId)
		return tmp.Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateByOrgAndAgentAndHashAndProcId 通过通过企业id\终端id\hash\ProcId修改
func (m *ProcessRunningModel) UpdateByOrgAndAgentAndHashAndProcId(orgId int64, agentId, hash string, procId int, updates map[string]interface{}) (bool, int64) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		tmp = tmp.Where("org_id=? AND agent_id=? AND proc_hash=? AND proc_id=?", orgId, agentId, hash, procId)
		return tmp.Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// GetItemByOrgAndAgentAndHashAndProcId 通过企业id\终端id\hash\ProcId获取一条记录
func (m *ProcessRunningModel) GetItemByOrgAndAgentAndHashAndProcId(orgId int64, agentId, hash string, procId int) (item *ProcessRunningModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND agent_id=? AND proc_hash=? AND proc_id=?", orgId, agentId, hash, procId).Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return nil
	}
}

// DeleteByConditions 通过条件删除
func (m *ProcessRunningModel) DeleteByConditions(condition SearchConditions) bool {
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
func (m *ProcessRunningModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*ProcessAndRunning) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {

		t1As := m.TableName() + "."
		t2Name := CreateProcessFactory().TableName()
		t2As := "t2."

		columns := []string{
			t1As + "org_id",
			t1As + "agent_id",
			t1As + "proc_hash",
			t1As + "proc_id",
			t1As + "id",
			t1As + "proc_dir",
			t1As + "proc_user",
			t1As + "created_at",
			t2As + "status",
			t2As + "sync_status_at",
			t2As + "self_status",
			t2As + "threaten_status",
			t2As + "auto_white_switch",
			t2As + "introduce",
			t2As + "copyright",
			t2As + "signature",
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
func (m *ProcessRunningModel) Count(condition SearchConditions) int64 {
	var total int64
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		t1As := m.TableName() + "."
		t2Name := CreateProcessFactory().TableName()
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
func (m *ProcessRunningModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*ProcessAndRunning) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
