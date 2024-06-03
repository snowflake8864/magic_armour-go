package proc_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
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
