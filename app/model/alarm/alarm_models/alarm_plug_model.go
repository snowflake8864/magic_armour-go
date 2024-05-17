package alarm_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// CreateAlarmPlugFactory 告警: 封堵
func CreateAlarmPlugFactory() *AlarmPlugModel {
	return &AlarmPlugModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAlarm)}}
}

//结构体

type AlarmPlugModel struct {
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
	Ip           string `json:"ip"`
	PlugType     int    `json:"plug_type"`
	Direction    int    `json:"direction"`
	Duration     int    `json:"duration"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time" gorm:"default:NULL"`
}

// 表名
func (m *AlarmPlugModel) TableName() string {
	return "a_alarm_plug"
}

func (m *AlarmPlugModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			ip VARCHAR(128) NOT NULL DEFAULT '' COMMENT '封堵ip',
			created_at DATETIME NULL,
			id BIGINT NOT NULL DEFAULT '0',
			type SMALLINT NOT NULL DEFAULT '0' COMMENT '类型',
			level TINYINT NOT NULL DEFAULT '0' COMMENT '级别',
			agent_at DATETIME NULL COMMENT '终端时间',
			handle_status TINYINT NOT NULL DEFAULT '0' COMMENT '处理状态(0未处理，1已处理)',
			handle_remark VARCHAR(512) NOT NULL DEFAULT '' COMMENT '处理备注',
			handle_user VARCHAR(128) NOT NULL DEFAULT '' COMMENT '处理人(user.username)',
			handle_at DATETIME NULL COMMENT '处理时间',
			plug_type SMALLINT NOT NULL DEFAULT '0' COMMENT '封堵类型',
			direction SMALLINT NOT NULL DEFAULT '0' COMMENT '方向',
			duration INT NOT NULL DEFAULT '0' COMMENT '封堵时长(分钟),0为永久',
			start_time DATETIME NULL COMMENT '封堵开始时间',
			end_time DATETIME NULL COMMENT '封堵结束时间',
			INDEX index_id (id) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id,ip,created_at)
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

// InsertItem 插入一条
func (m *AlarmPlugModel) InsertItem(s *AlarmPlugModel) bool {

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
func (m *AlarmPlugModel) InsertBatch(list *[]AlarmPlugModel) bool {
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
func (m *AlarmPlugModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) bool {
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
func (m *AlarmPlugModel) DeleteByConditions(condition SearchConditions) bool {
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
func (m *AlarmPlugModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*AlarmPlugModel) {
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
func (m *AlarmPlugModel) Count(condition SearchConditions) int64 {
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
func (m *AlarmPlugModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*AlarmPlugModel) {
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
