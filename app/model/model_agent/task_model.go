package model_agent

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

//  终端任务表

func CreateTaskFactory() *TaskModel {
	return &TaskModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type TaskModel struct {
	model.BaseModel
	Id          string `json:"id" gorm:"primaryKey"`
	OrgId       int64  `json:"org_id"`
	AgentId     string `json:"agent_id"`
	Type        int    `json:"type"`
	Remark      string `json:"remark"`
	Status      int    `json:"status"`
	NetIp       string `json:"net_ip"`
	CreatedAt   string `json:"created_at"`
	CompletedAt string `json:"completed_at"`
}

// 表名

func (m *TaskModel) TableName() string {
	return "a_tasks"
}

func (m *TaskModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			type smallint  NOT NULL DEFAULT '0' COMMENT '类型',
			created_at datetime NULL,
 			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'md5值',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '备注',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(0进行中, 1已完成)',
			net_ip varchar(128)  NOT NULL DEFAULT '' COMMENT '通讯ip',
			completed_at datetime NULL COMMENT '任务完成时间'
		)
		UNIQUE KEY(org_id, agent_id, type, created_at)
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

// InsertBatch 插入多条
func (m *TaskModel) InsertBatch(list *[]TaskModel) bool {
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
func (m *TaskModel) UpdateByConditions(condition map[string]interface{}, updates map[string]interface{}) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
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

// DeleteByOrg 通过企业i删除
func (m *TaskModel) DeleteByOrg(orgId int64) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *TaskModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*TaskModel, err error) {
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
func (m *TaskModel) Count(condition map[string]interface{}) int64 {
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
func (m *TaskModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*TaskModel) {
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
