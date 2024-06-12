package port_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 表: 终端端口虚开策略

func CreateVirtualPolicyFactory() *VirtualPolicyModel {
	return &VirtualPolicyModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type VirtualPolicyModel struct {
	basePort
	model.BaseModel
	Id          int64  `json:"id" gorm:"primaryKey"`
	OrgId       int64  `json:"org_id"`
	AgentId     string `json:"agent_id"`
	Type        string `json:"type"`
	Protocol    string `json:"protocol"`
	SourceIp    string `json:"source_ip"`
	SourcePort  string `json:"source_port"`
	ForwardIp   string `json:"forward_ip"`
	ForwardPort string `json:"forward_port"`
	AlarmLevel  int    `json:"alarm_level"`
	OpenSwitch  int    `json:"open_switch"`
	CreateId    int64  `json:"create_id"`
	CreateName  string `json:"create_name"`
	CreatedAt   string `json:"created_at"`
	UpdateId    int64  `json:"update_id"`
	UpdateName  string `json:"update_name"`
	UpdatedAt   string `json:"updated_at"`
}

// 表名

func (m *VirtualPolicyModel) TableName() string {
	return "a_port_virtual_policy"
}

func (m *VirtualPolicyModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			id BIGINT NOT NULL DEFAULT '0',
			type varchar(8)  NOT NULL DEFAULT '' COMMENT 'ip类型(ipv4,ipv6)',
			protocol varchar(8)  NOT NULL DEFAULT '' COMMENT '传输协议(tcp,udp)',
			source_ip varchar(128)  NOT NULL DEFAULT '' COMMENT '源ip',
			source_port varchar(512)  NOT NULL DEFAULT '' COMMENT '源端口:单端口(80);多端口(80,90);范围(80-90)',
			forward_ip varchar(128)  NOT NULL DEFAULT '' COMMENT '转发ip',
			forward_port varchar(64)  NOT NULL DEFAULT '' COMMENT '转发端口',
			alarm_level tinyint NOT NULL DEFAULT '0' COMMENT '告警级别',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			update_id int NOT NULL DEFAULT '0' COMMENT '上一次编辑用户(user.id)',
			update_name varchar(128)  NOT NULL DEFAULT '' COMMENT '上一次编辑用户(user.username)',
			updated_at datetime NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id, id) 
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

func (m *VirtualPolicyModel) InsertItem(s *VirtualPolicyModel) bool {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(s)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// UpdateByConditions 通过条件修改
func (m *VirtualPolicyModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) (bool, int64) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		return tmp.Updates(updates)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// DeleteByConditions 通过条件删除
func (m *VirtualPolicyModel) DeleteByConditions(condition SearchConditions) bool {
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

// GetItemByConditions 通过id获取一条记录
func (m *VirtualPolicyModel) GetItemByConditions(condition SearchConditions) (item *VirtualPolicyModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
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
		return nil
	}
}

// GetListByConditions 获取符合条件的全部记录
func (m *VirtualPolicyModel) GetListByConditions(condition SearchConditions) (items []*VirtualPolicyModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *VirtualPolicyModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*VirtualPolicyModel) {
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
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return nil
	}
}

// Count 通过条件统计
func (m *VirtualPolicyModel) Count(condition SearchConditions) int64 {
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
func (m *VirtualPolicyModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*VirtualPolicyModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}

// CalcByConditionsGroupAgent 通过条件获取多条记录且按终端分组
func (m *VirtualPolicyModel) CalcByConditionsGroupAgent(condition SearchConditions, column string) (items []*Calc) {

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
func (m *VirtualPolicyModel) CalcByConditionsGroupAgentReturnMap(condition SearchConditions, column string) map[string]interface{} {
	var result = make(map[string]interface{})

	items := m.CalcByConditionsGroupAgent(condition, column)
	if items == nil || len(items) == 0 {
		return result
	}
	for _, item := range items {
		result[item.AgentId] = *item
	}

	return result
}
