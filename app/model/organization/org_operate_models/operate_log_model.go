package org_operate_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 角色表

func CreateLogFactory() *LogModel {
	return &LogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type LogModel struct {
	baseOperate
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	Path       string `json:"path"`
	Module     string `json:"module"`
	NetIp      string `json:"net_ip"`
	Input      string `json:"input"`
	AgentId    string `json:"agent_id"`
	Method     string `json:"method"`
	Status     int    `json:"status"`
	HttpCode   int    `json:"http_code"`
	DelayTime  int64  `json:"delay_time"`
	Resp       string `json:"resp"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
}

// 表名
func (m *LogModel) TableName() string {
	return "o_operate_logs"
}

func (m *LogModel) CreateTable1() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(	
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			module varchar(32) NOT NULL DEFAULT '' COMMENT '业务模块',
			created_at datetime NOT NULL,
			path varchar(128)  NOT NULL DEFAULT '' COMMENT '路径',
			net_ip varchar(128) NOT NULL DEFAULT '' COMMENT '通讯ip',
			id BIGINT NOT NULL DEFAULT '0',
			input text NULL COMMENT '请求参数',
			agent_id varchar(128) NOT NULL DEFAULT ''COMMENT '终端id(agent.id)',
			method varchar(8) NULL COMMENT '请求方法',
			status tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态(1成功 0失败)',
			http_code smallint(4) NOT NULL DEFAULT '0' COMMENT 'http状态码',
			delay_time int NOT NULL DEFAULT '0' COMMENT '延迟时间(ms)',
			resp text NULL COMMENT '返回结果',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)'
		)
		DUPLICATE  KEY(org_id, module, created_at)
		DISTRIBUTED  BY HASH(org_id) BUCKETS AUTO 
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
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *LogModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*LogModel, err error) {
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
		return items, nil
	} else {
		return nil, result.Error
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
