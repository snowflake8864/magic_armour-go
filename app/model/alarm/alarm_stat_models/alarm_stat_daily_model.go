package alarm_stat_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// CreateAlarmStatDailyFactory 告警统计: 每天的
func CreateAlarmStatDailyFactory() *AlarmStatDailyModel {
	return &AlarmStatDailyModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAlarm)}}
}

//结构体

type AlarmStatDailyModel struct {
	baseAlarm
	model.BaseModel
	OrgId               int64  `json:"org_id"`
	AgentId             string `json:"agent_id"`
	StatDay             string `json:"stat_day"`
	StatMonth           int    `json:"stat_month"`
	Total               int    `json:"total"`
	Handled             int    `json:"handled"`
	AddTotal            int    `json:"add_total"`
	AddHandled          int    `json:"add_handled"`
	ProcessTotal        int    `json:"process_total"`
	ProcessHandled      int    `json:"process_handled"`
	ProcessAddTotal     int    `json:"process_add_total"`
	ProcessAddHandled   int    `json:"process_add_handled"`
	DllTotal            int    `json:"dll_total"`
	DllHandled          int    `json:"dll_handled"`
	DllAddTotal         int    `json:"dll_add_total"`
	DllAddHandled       int    `json:"dll_add_handled"`
	BlackmailTotal      int    `json:"blackmail_total"`
	BlackmailHandled    int    `json:"blackmail_handled"`
	BlackmailAddTotal   int    `json:"blackmail_add_total"`
	BlackmailAddHandled int    `json:"blackmail_add_handled"`
	TamperTotal         int    `json:"tamper_total"`
	TamperHandled       int    `json:"tamper_handled"`
	TamperAddTotal      int    `json:"tamper_add_total"`
	TamperAddHandled    int    `json:"tamper_add_handled"`
	PortTotal           int    `json:"port_total"`
	PortHandled         int    `json:"port_handled"`
	PortAddTotal        int    `json:"port_add_total"`
	PortAddHandled      int    `json:"port_add_handled"`
	PlugTotal           int    `json:"plug_total"`
	PlugHandled         int    `json:"plug_handled"`
	PlugAddTotal        int    `json:"plug_add_total"`
	PlugAddHandled      int    `json:"plug_add_handled"`
	UsbTotal            int    `json:"usb_total"`
	UsbHandled          int    `json:"usb_handled"`
	UsbAddTotal         int    `json:"usb_add_total"`
	UsbAddHandled       int    `json:"usb_add_handled"`
	OutreachTotal       int    `json:"outreach_total"`
	OutreachHandled     int    `json:"outreach_handled"`
	OutreachAddTotal    int    `json:"outreach_add_total"`
	OutreachAddHandled  int    `json:"outreach_add_handled"`
	LoginTotal          int    `json:"login_total"`
	LoginHandled        int    `json:"login_handled"`
	LoginAddTotal       int    `json:"login_add_total"`
	LoginAddHandled     int    `json:"login_add_handled"`
	SelfTotal           int    `json:"self_total"`
	SelfHandled         int    `json:"self_handled"`
	SelfAddTotal        int    `json:"self_add_total"`
	SelfAddHandled      int    `json:"self_add_handled"`
	CreatedAt           string `json:"created_at"`
}

// 表名
func (m *AlarmStatDailyModel) TableName() string {
	return "a_alarm_stat_daily"
}

func (m *AlarmStatDailyModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
            stat_month int(11) NOT NULL COMMENT '统计：日期',
			stat_day DATEV2 NOT NULL COMMENT '统计：日期',
			total int(11) NOT NULL DEFAULT '0' COMMENT '告警:总数',
			handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:已处理数',
			add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:进新增-总数',
			add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:新增-已处理数',
			process_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:进程-总数',
			process_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:进程-已处理数',
			process_add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:进程-新增-总数',
			process_add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:进程-新增-已处理数',
			dll_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:模块-总数',
			dll_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:模块-已处理数',
			dll_add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:模块-新增-总数',
			dll_add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:模块-新增-已处理数',
			blackmail_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:勒索-总数',
			blackmail_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:勒索-已处理数',
			blackmail_add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:勒索-新增-总数',
			blackmail_add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:勒索-新增-已处理数',
			tamper_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:勒索-总数',
			tamper_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:勒索-已处理数',
			tamper_add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:勒索-新增-总数',
			tamper_add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:勒索-新增-已处理数',
			port_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:端口-总数',
			port_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:端口-已处理数',
			port_add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:端口-新增-总数',
			port_add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:端口-新增-已处理数',
			plug_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:封堵-总数',
			plug_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:封堵-已处理数',
			plug_add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:封堵-新增-总数',
			plug_add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:封堵-新增-已处理数',
			usb_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:外设-总数',
			usb_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:外设-已处理数',
			usb_add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:外设-新增-总数',
			usb_add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:外设-新增-已处理数',
			outreach_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:外联-总数',
			outreach_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:外联-已处理数',
			outreach_add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:外联-新增-总数',
			outreach_add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:外联-新增-已处理数',
			login_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:登录破解-总数',
			login_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:登录破解-已处理数',
			login_add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:登录破解-新增-总数',
			login_add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:登录破解-新增-已处理数',
			self_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:自保-总数',
			self_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:自保-已处理数',
			self_add_total int(11) NOT NULL DEFAULT '0' COMMENT '告警:自保-新增-总数',
			self_add_handled int(11) NOT NULL DEFAULT '0' COMMENT '告警:自保-新增-已处理数',
			created_at datetime NULL
		)
		UNIQUE KEY(org_id, agent_id,stat_month,stat_day)
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
func (m *AlarmStatDailyModel) InsertItem(s *AlarmStatDailyModel) bool {

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
func (m *AlarmStatDailyModel) InsertBatch(list *[]AlarmStatDailyModel) bool {
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
func (m *AlarmStatDailyModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !helpers.StringIsEmpty(whereSQL) {
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
func (m *AlarmStatDailyModel) DeleteByConditions(condition SearchConditions) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !helpers.StringIsEmpty(whereSQL) {
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
func (m *AlarmStatDailyModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*AlarmStatDailyModel) {
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
		return items
	} else {
		return nil
	}
}

// CalcByConditions 通过条件获取多条记录
func (m *AlarmStatDailyModel) CalcByConditions(condition SearchConditions, column string) (item *AlarmStatDailyModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition))
		tmp := tx.Model(m).Select(column).Where(whereSQL, vals...)
		return tmp.Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item
	} else {
		return nil
	}
}

// CalcByConditionsGroupDay 通过条件获取多条记录且按天分组
func (m *AlarmStatDailyModel) CalcByConditionsGroupDay(condition SearchConditions, column string) (items []*AlarmStatDailyModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition))
		tmp := tx.Model(m).Select(column).Where(whereSQL, vals...)
		return tmp.Group("stat_day").Order("stat_day ASC").Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items
	} else {
		return nil
	}
}
