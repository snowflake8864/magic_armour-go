package proc_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
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

// UpdateByOrgAndAgentAndHashAndProcId 通过通过企业id\终端id\hash\ProcId修改
func (m *DllRunningModel) UpdateByOrgAndAgentAndHashAndProcId(orgId int64, agentId, hash string, procId int, updates map[string]interface{}) (bool, int64) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		tmp = tmp.Where("org_id=? AND agent_id=? AND dll_hash=? AND proc_id", orgId, agentId, hash, procId)
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
func (m *DllRunningModel) GetItemByOrgAndAgentAndHashAndProcId(orgId int64, agentId, hash string, procId int) (item *DllRunningModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND agent_id=? AND dll_hash=? AND proc_id", orgId, agentId, hash, procId).Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return nil
	}
}
