package model_alarm

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// CreateAlarmSelfFactory 告警: 自保
func CreateAlarmSelfFactory() *AlarmSelfModel {
	return &AlarmSelfModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAlarm)}}
}

//结构体

type AlarmSelfModel struct {
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
	ProcParam    string `json:"proc_param"`
	FileDir      string `json:"file_dir"`
	TargetDir    string `json:"target_dir"`
}

// 表名
func (m *AlarmSelfModel) TableName() string {
	return "a_alarm_self"
}

func (m *AlarmSelfModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			proc_hash VARCHAR(64) NOT NULL DEFAULT '' COMMENT '进程:hash',
			type SMALLINT NOT NULL DEFAULT '0' COMMENT '类型',
			created_at DATETIME NULL,
			agent_at DATETIME NULL COMMENT '终端时间',
			id BIGINT NOT NULL DEFAULT '0',
			level TINYINT NOT NULL DEFAULT '0' COMMENT '级别',
			handle_status TINYINT NOT NULL DEFAULT '0' COMMENT '处理状态(0未处理，1已处理)',
			handle_remark VARCHAR(512) NOT NULL DEFAULT '' COMMENT '处理备注',
			handle_user VARCHAR(128) NOT NULL DEFAULT '' COMMENT '处理人(user.username)',
			handle_at DATETIME NULL COMMENT '处理时间',
			proc_dir VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '进程:路径',
			proc_param text NULL COMMENT '进程:参数',
            file_dir VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '文件:路径',
            target_dir VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '目标:路径',
			INDEX index_id (id) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED,
			INDEX index_proc_dir (proc_dir) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id,proc_hash,type,created_at,agent_at)
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
func (m *AlarmSelfModel) InsertItem(s *AlarmSelfModel) bool {

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
func (m *AlarmSelfModel) InsertBatch(list *[]AlarmSelfModel) bool {
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
func (m *AlarmSelfModel) UpdateByConditions(condition map[string]interface{}, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where(condition).Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *AlarmSelfModel) DeleteByConditions(condition map[string]interface{}) bool {
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
func (m *AlarmSelfModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*AlarmSelfModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := m.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals)
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
func (m *AlarmSelfModel) Count(condition map[string]interface{}) int64 {
	var total int64
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := m.Model(m).Select("id")
		if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals)
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
func (m *AlarmSelfModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*AlarmSelfModel) {
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
