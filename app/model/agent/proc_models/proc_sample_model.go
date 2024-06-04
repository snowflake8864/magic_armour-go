package proc_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 表: 进程与模块样本

func CreateProcSampleFactory() *ProcSampleModel {
	return &ProcSampleModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ProcSampleModel struct {
	baseProc
	model.BaseModel
	Id           string `json:"id" gorm:"primaryKey"`
	OrgId        int64  `json:"org_id"`
	AgentId      string `json:"agent_id"`
	Dir          string `json:"dir"`
	Introduce    string `json:"introduce"`
	Copyright    string `json:"copyright"`
	Signature    string `json:"signature"`
	Attribute    string `json:"attribute"`
	DataFrom     int    `json:"data_from"`
	Status       int    `json:"status"`
	OriginalName string `json:"original_name"`
	StoragePath  string `json:"storage_path"`
	FullPath     string `json:"full_path"`
	CreateId     int64  `json:"create_id"`
	CreateName   string `json:"create_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// 表名

func (m *ProcSampleModel) TableName() string {
	return "a_proc_samples"
}

func (m *ProcSampleModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'hash值',
			dir varchar(512)  NOT NULL DEFAULT '' COMMENT '路径',
			introduce varchar(256)  NOT NULL DEFAULT '' COMMENT '描述',
			copyright varchar(256)  NOT NULL DEFAULT '' COMMENT '版权',
			signature varchar(512)  NOT NULL DEFAULT '' COMMENT '签名',
			attribute text  NOT NULL DEFAULT '' COMMENT '属性',
			data_from tinyint NOT NULL DEFAULT '1' COMMENT '数据来源(1进程 2模块)',
			status tinyint NOT NULL DEFAULT '1' COMMENT '状态',
			original_name varchar(128)  NOT NULL DEFAULT '' COMMENT '原文件名',
			storage_path varchar(512)  NOT NULL DEFAULT '' COMMENT '存储相对路径(项目下)',
			full_path varchar(1024)  NOT NULL DEFAULT '' COMMENT '完整路径',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id,agent_id,id) 
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

func (m *ProcSampleModel) InsertItem(s *ProcSampleModel) bool {

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
func (m *ProcSampleModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) (bool, int64) {
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
func (m *ProcSampleModel) DeleteByConditions(condition SearchConditions) bool {
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

// GetItemByOrgAndAgentAndHash 通过企业id\终端id\hash\ProcId获取一条记录
func (m *ProcSampleModel) GetItemByOrgAndAgentAndHash(orgId int64, agentId, hash string) (item *ProcSampleModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND agent_id=? AND id=? AND proc_id=?", orgId, agentId, hash).Take(&item)
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
func (m *ProcSampleModel) GetListByConditions(condition SearchConditions) (items []*ProcSampleModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *ProcSampleModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*ProcSampleModel) {
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
func (m *ProcSampleModel) Count(condition SearchConditions) int64 {
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
func (m *ProcSampleModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*ProcSampleModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
