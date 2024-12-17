package plug_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 表: 终端封堵IP名单库

func CreateIpPolicyFactory() *IpPolicyModel {
	return &IpPolicyModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type IpPolicyModel struct {
	basePlug
	model.BaseModel
	Id            int64  `json:"id" gorm:"primaryKey"`
	OrgId         int64  `json:"org_id"`
	AgentId       string `json:"agent_id"`
	Ip            string `json:"ip"`
	Type          int    `json:"type"`
	Direction     int    `json:"direction"`
	Remark        string `json:"remark"`
	EffectMode    int    `json:"effect_mode"`
	EffectStartAt string `json:"effect_start_at"  gorm:"default:NULL"`
	EffectEndAt   string `json:"effect_end_at"  gorm:"default:NULL"`
	OpenSwitch    int    `json:"open_switch"`
	CreateId      int64  `json:"create_id"`
	CreateName    string `json:"create_name"`
	CreatedAt     string `json:"created_at"`
	UpdateId      int64  `json:"update_id"`
	UpdateName    string `json:"update_name"`
	UpdatedAt     string `json:"updated_at"`
}

// 表名

func (m *IpPolicyModel) TableName() string {
	return "a_plug_ip_policy"
}

func (m *IpPolicyModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			id BIGINT NOT NULL DEFAULT '0',
			type tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型(1白 2黑)',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			ip varchar(128)  NOT NULL DEFAULT '' COMMENT 'ip',
			direction tinyint(1) NOT NULL DEFAULT '0' COMMENT '方向(1IN 2OUT),二进制表示',
			remark varchar(256)  NOT NULL DEFAULT '' COMMENT '备注',
			effect_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '生效模式',
			effect_start_at dateV2 NULL COMMENT '生效开始时间',
			effect_end_at dateV2 NULL COMMENT '生效结束时间',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			update_id int NOT NULL DEFAULT '0' COMMENT '上一次编辑用户(user.id)',
			update_name varchar(128)  NOT NULL DEFAULT '' COMMENT '上一次编辑用户(user.username)',
			updated_at datetime NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id, agent_id, id, type,create_id) 
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

func (m *IpPolicyModel) InsertItem(s *IpPolicyModel) bool {

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
func (m *IpPolicyModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) (bool, int64) {
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
func (m *IpPolicyModel) DeleteByConditions(condition SearchConditions) bool {
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
func (m *IpPolicyModel) GetItemByConditions(condition SearchConditions) (item *IpPolicyModel) {
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
func (m *IpPolicyModel) GetListByConditions(condition SearchConditions) (items []*IpPolicyModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *IpPolicyModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*IpPolicyModel) {
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
func (m *IpPolicyModel) Count(condition SearchConditions) int64 {
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
func (m *IpPolicyModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*IpPolicyModel) {
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

// CalcByConditionsGroupAgent 通过条件获取多条记录且按终端分组
func (m *IpPolicyModel) CalcByConditionsGroupAgent(condition SearchConditions, column string) (items []*Calc) {

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
func (m *IpPolicyModel) CalcByConditionsGroupAgentReturnMap(condition SearchConditions, column string) map[string]*Calc {
	var result = make(map[string]*Calc)

	items := m.CalcByConditionsGroupAgent(condition, column)
	if items == nil || len(items) == 0 {
		return result
	}
	for _, item := range items {
		result[item.AgentId] = item
	}

	return result
}
