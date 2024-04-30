package model_alarm

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// CreateAlarmPortFactory 告警: 端口
func CreateAlarmPortFactory() *AlarmPortModel {
	return &AlarmPortModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAlarm)}}
}

//结构体

type AlarmPortModel struct {
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
	Weight       int    `json:"weight"`
	AttackIp     string `json:"attack_ip"`
	TargetIp     string `json:"target_ip"`
	TargetPort   int    `json:"target_port"`
	RedirectIp   string `json:"redirect_ip"`
	RedirectPort int    `json:"redirect_port"`
	AttackTotal  int    `json:"attack_total"`
}

// 表名
func (m *AlarmPortModel) TableName() string {
	return "a_alarm_port"
}
func (m *AlarmPortModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			attack_ip VARCHAR(128) NOT NULL DEFAULT '' COMMENT '攻击ip',
			target_ip VARCHAR(128) NOT NULL DEFAULT '' COMMENT '目标ip',
			target_port INT NOT NULL DEFAULT '0' COMMENT '目标端口',	
			created_at DATETIME NULL,
			id BIGINT NOT NULL DEFAULT '0',
			type SMALLINT NOT NULL DEFAULT '0' COMMENT '类型',
			level TINYINT NOT NULL DEFAULT '0' COMMENT '级别',
			agent_at DATETIME NULL COMMENT '终端时间',
			handle_status TINYINT NOT NULL DEFAULT '0' COMMENT '处理状态(0未处理，1已处理)',
			handle_remark VARCHAR(512) NOT NULL DEFAULT '' COMMENT '处理备注',
			handle_user VARCHAR(128) NOT NULL DEFAULT '' COMMENT '处理人(user.username)',
			handle_at DATETIME NULL COMMENT '处理时间',
			weight SMALLINT NOT NULL DEFAULT '0' COMMENT '权重',
			redirect_ip VARCHAR(128) NOT NULL DEFAULT '' COMMENT '重定向ip',
			redirect_port INT NOT NULL DEFAULT '0' COMMENT '重定向端口',
			attack_total INT NOT NULL DEFAULT '0' COMMENT '攻击次数',
			INDEX index_id (id) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id,attack_ip,target_ip,target_port,created_at)
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
func (m *AlarmPortModel) InsertItem(s *AlarmPortModel) bool {

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
func (m *AlarmPortModel) InsertBatch(list *[]AlarmPortModel) bool {
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
func (m *AlarmPortModel) UpdateByConditions(condition map[string]interface{}, updates map[string]interface{}) bool {
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
func (m *AlarmPortModel) DeleteByConditions(condition map[string]interface{}) bool {
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
func (m *AlarmPortModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*AlarmPortModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := m.Model(m)
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
func (m *AlarmPortModel) Count(condition map[string]interface{}) int64 {
	var total int64
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := m.Model(m).Select("id")
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
func (m *AlarmPortModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*AlarmPortModel) {
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
