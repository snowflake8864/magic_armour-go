package usb_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 外设插拔日志表

func CreateLogFactory() *LogModel {
	return &LogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type LogModel struct {
	baseUsb
	model.BaseModel
	OrgId      int64  `json:"org_id"`
	AgentId    string `json:"agent_id"`
	UsbHash    string `json:"usb_hash"`
	Name       string `json:"name"`
	Status     int    `json:"status"`
	PlugStatus int    `json:"plug_status"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *LogModel) TableName() string {
	return "a_usb_logs"
}

func (m *LogModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			usb_hash varchar(64)  NOT NULL DEFAULT '' COMMENT '设备:hash',
			created_at datetime NULL,
			name varchar(512)  NOT NULL DEFAULT '' COMMENT '设备:名称',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(冗余)',
			plug_status tinyint NOT NULL DEFAULT '1' COMMENT '插拔状态(1拔出 2插入)'
		)
		DUPLICATE KEY(org_id,agent_id, usb_hash, created_at)
		DISTRIBUTED BY HASH(org_id) BUCKETS AUTO 
		PROPERTIES (
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

func (m *LogModel) InsertItem(s *LogModel) bool {

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

// InsertBatch 插入多条
func (m *LogModel) InsertBatch(list *[]LogModel) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(list)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *LogModel) DeleteByConditions(condition SearchConditions) bool {
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
func (m *LogModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*LogModel) {
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
func (m *LogModel) Count(condition SearchConditions) int64 {
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
func (m *LogModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*LogModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
