package model_syslog

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// CreateSyslogInnerFactory 通讯日志：对内通讯
func CreateSyslogInnerFactory() *SyslogInnerModel {
	return &SyslogInnerModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAlarm)}}
}

// 结构体
type SyslogInnerModel struct {
	model.BaseModel
	Id               int64  `json:"id" gorm:"primaryKey" `
	OrgId            int64  `json:"org_id"`
	AgentId          string `json:"agent_id"`
	AgentAt          string `json:"agent_at"`
	CreatedAt        string `json:"created_at"`
	ProcHash         string `json:"proc_hash"`
	ProcId           int    `json:"proc_id"`
	ProcDir          string `json:"proc_dir"`
	TargetIp         string `json:"target_ip"`
	TargetPort       int    `json:"target_port"`
	TargetPortStatus int    `json:"target_port_status"`
	SourceIp         string `json:"source_ip"`
	SourcePort       int    `json:"source_port"`
	Proto            int    `json:"proto"`
}

// 表名
func (m *SyslogInnerModel) TableName() string {
	return "a_syslog_inner"
}

func (m *SyslogInnerModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			created_at DATETIME NULL,
			proc_hash VARCHAR(64) NOT NULL DEFAULT '' COMMENT '进程:hash',
			target_ip VARCHAR(128) NOT NULL DEFAULT '' COMMENT '目标ip',
			target_port INT NOT NULL DEFAULT '0' COMMENT '目标端口',
			source_ip VARCHAR(128) NOT NULL DEFAULT '' COMMENT '源ip',
			source_port INT NOT NULL DEFAULT '0' COMMENT '源端口',	
			id BIGINT NOT NULL DEFAULT '0',
			agent_at DATETIME NULL COMMENT '终端时间',
			proc_id INT NOT NULL DEFAULT '0' COMMENT '进程:ID',
			proc_dir VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '进程:路径',
			target_port_status TINYINT NOT NULL DEFAULT '0' COMMENT '目标端口状态',
			proto TINYINT NOT NULL DEFAULT '0' COMMENT '协议'
		)
		DUPLICATE KEY(org_id,agent_id,created_at,proc_hash,target_ip,target_port,source_ip,source_port,id) 
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
func (m *SyslogInnerModel) InsertItem(s *SyslogInnerModel) bool {

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
func (m *SyslogInnerModel) InsertBatch(list *[]SyslogInnerModel) bool {
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
func (m *SyslogInnerModel) DeleteByConditions(condition map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
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
func (m *SyslogInnerModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*SyslogInnerModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
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
func (m *SyslogInnerModel) Count(condition map[string]interface{}) int64 {
	var total int64
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Select("id")
		if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
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
func (m *SyslogInnerModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*SyslogInnerModel) {
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
