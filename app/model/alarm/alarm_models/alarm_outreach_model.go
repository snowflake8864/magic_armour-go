package alarm_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// CreateOutreachFactory 告警: 外联
func CreateOutreachFactory() *OutreachModel {
	return &OutreachModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAlarm)}}
}

//结构体

type OutreachModel struct {
	baseAlarm
	model.BaseModel
	Id           int64  `json:"id" gorm:"primaryKey"`
	OrgId        int64  `json:"org_id"`
	AgentId      string `json:"agent_id"`
	Type         int    `json:"type"`
	Level        int    `json:"level"`
	AgentAt      string `json:"agent_at"`
	HandleStatus int    `json:"handle_status"`
	HandleRemark string `json:"handle_remark"`
	HandleUser   string `json:"handle_user"`
	HandleAt     string `json:"handle_at"`
	CreatedAt    string `json:"created_at"`
	ProcHash     string `json:"proc_hash"`
	ProcDir      string `json:"proc_dir"`
	ProcId       int    `json:"proc_id"`
	DomainName   string `json:"domain_name"`
	TargetIp     string `json:"target_ip"`
}

// 表名
func (m *OutreachModel) TableName() string {
	return "a_alarm_outreach"
}

func (m *OutreachModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			proc_hash VARCHAR(64) NOT NULL DEFAULT '' COMMENT '进程:hash',
			domain_name VARCHAR(512) NOT NULL DEFAULT '' COMMENT '域名',
			target_ip VARCHAR(512) NOT NULL DEFAULT '' COMMENT '目标ip',
			created_at DATETIME NULL,
			agent_at DATETIME NULL COMMENT '终端时间',
			id BIGINT NOT NULL DEFAULT '0',
			type SMALLINT NOT NULL DEFAULT '0' COMMENT '类型',
			level TINYINT NOT NULL DEFAULT '0' COMMENT '级别',
			handle_status TINYINT NOT NULL DEFAULT '0' COMMENT '处理状态(0未处理，1已处理)',
			handle_remark VARCHAR(512) NOT NULL DEFAULT '' COMMENT '处理备注',
			handle_user VARCHAR(128) NOT NULL DEFAULT '' COMMENT '处理人(user.username)',
			handle_at DATETIME NULL COMMENT '处理时间',
			proc_dir VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '进程:路径',
			proc_id int NOT NULL DEFAULT '0' COMMENT '进程:id',
			INDEX index_id (id) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id,proc_hash,domain_name,target_ip,created_at,agent_at)
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

func (m *OutreachModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *OutreachModel) InsertItem(s *OutreachModel) bool {

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
func (m *OutreachModel) InsertBatch(list *[]OutreachModel) bool {
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

// UpdateByConditions 通过条件修改
func (m *OutreachModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *OutreachModel) DeleteByConditions(condition SearchConditions) bool {
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
func (m *OutreachModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*OutreachModel) {
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
		return nil
	}
}

// Count 通过条件统计
func (m *OutreachModel) Count(condition SearchConditions) int64 {
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
func (m *OutreachModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*OutreachModel) {
	total = m.Count(condition)
	if total == 0 {
		return total, nil
	}
	items = m.GetItemsByConditions(condition, offset, limit, order)
	if len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}

// CalcByConditions 通过条件获取多条记录且按终端分组
func (m *OutreachModel) CalcByConditions(condition SearchConditions, column string) (item *AlarmCalc) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m).Select(column)
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
		return item
	}
}

// CalcByConditionsGroupAgent 通过条件获取多条记录且按终端分组
func (m *OutreachModel) CalcByConditionsGroupAgent(condition SearchConditions, column string) (items []*AlarmCalc) {

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
func (m *OutreachModel) CalcByConditionsGroupAgentReturnMap(condition SearchConditions, column string) map[string]*AlarmCalc {
	var result = make(map[string]*AlarmCalc)

	items := m.CalcByConditionsGroupAgent(condition, column)
	if items == nil || len(items) == 0 {
		return result
	}
	for _, item := range items {
		result[item.AgentId] = item
	}

	return result
}
