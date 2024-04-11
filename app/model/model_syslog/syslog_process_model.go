package model_syslog

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// CreateSyslogProcessFactory 通讯日志：进程启动
func CreateSyslogProcessFactory() *SyslogProcessModel {
	return &SyslogProcessModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAlarm)}}
}

// 结构体
type SyslogProcessModel struct {
	model.BaseModel
	Id          int64  `json:"id" gorm:"primaryKey" `
	OrgId       int64  `json:"org_id"`
	AgentId     string `json:"agent_id"`
	AgentAt     string `json:"agent_at"`
	CreatedAt   string `json:"created_at"`
	ProcHash    string `json:"proc_hash"`
	ProcId      int    `json:"proc_id"`
	ProcDir     string `json:"proc_dir"`
	ProcParam   string `json:"proc_param"`
	ParentHash  string `json:"parent_hash"`
	ParentId    int    `json:"parent_id"`
	ParentDir   string `json:"parent_dir"`
	ParentParam string `json:"parent_param"`
}

// 表名

func (m *SyslogProcessModel) TableName() string {
	return "a_syslog_process"
}

func (m *SyslogProcessModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			created_at DATETIME NULL,
			proc_hash VARCHAR(64) NOT NULL DEFAULT '' COMMENT '进程:hash',
			proc_id INT NOT NULL DEFAULT '0' COMMENT '进程:ID',
			parent_hash VARCHAR(64) NOT NULL DEFAULT '' COMMENT '父进程:hash',
			parent_id INT NOT NULL DEFAULT '0' COMMENT '父进程:ID',
			id BIGINT NOT NULL DEFAULT '0',
			agent_at DATETIME NULL COMMENT '终端时间',
			proc_dir VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '进程:路径',
			proc_param text NULL COMMENT '进程:参数',
			parent_dir VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '父进程:路径',
			parent_param text NULL COMMENT '父进程:参数'
		)
		DUPLICATE KEY(org_id,agent_id,created_at,proc_hash,proc_id,parent_hash,parent_id,id) 
		PARTITION BY RANGE(created_at)() 
		DISTRIBUTED  BY HASH(org_id) BUCKETS AUTO 
		PROPERTIES (
			'replication_allocation' = 'tag.location.default: 1',
			'dynamic_partition.enable' = 'true',
			"dynamic_partition.prefix" = "p",
			'dynamic_partition.time_unit' = 'MONTH',
			'dynamic_partition.start' = '-12',
			 "dynamic_partition.end" = "2",
			'dynamic_partition.create_history_partition' = 'true',
			'dynamic_partition.history_partition_num' = '-1'
		);
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *SyslogProcessModel) InsertItem(s *SyslogProcessModel) bool {

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

// InsertBatch 插入多条
func (m *SyslogProcessModel) InsertBatch(list *[]SyslogProcessModel) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(list)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *SyslogProcessModel) DeleteByConditions(condition map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where(condition).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *SyslogProcessModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*SyslogProcessModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Where(condition)
		if limit > 0 {
			tmp = tmp.Offset(offset).Limit(limit)
		}
		if !helpers.StringIsEmpty(order) {
			tmp = tmp.Order(order)
		} else {
			tmp = tmp.Order("created_at ASC")
		}
		return tmp.Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// Count 通过条件统计
func (m *SyslogProcessModel) Count(condition map[string]interface{}) int64 {
	var total int64
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Select("id").Where(condition).Count(&total)
	})
	result := m.Raw(m.TrimSql(sql)).Count(&total)
	if result.Error == nil {
		return total
	} else {
		return 0
	}
}

// Page 通过条件分页
func (m *SyslogProcessModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*SyslogProcessModel) {
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
