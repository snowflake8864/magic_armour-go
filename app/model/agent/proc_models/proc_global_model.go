package proc_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 表: 进程与模块全局名单

func CreateProcGlobalFactory() *ProcGlobalModel {
	return &ProcGlobalModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type ProcGlobalModel struct {
	model.BaseModel
	Id         string `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	Type       int    `json:"type"`
	Dir        string `json:"dir"`
	Remark     string `json:"remark"`
	DataFrom   int    `json:"data_from"`
	OpenSwitch int    `json:"open_switch"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *ProcGlobalModel) TableName() string {
	return "a_proc_global"
}

func (m *ProcGlobalModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			id varchar(64)  NOT NULL DEFAULT '' COMMENT 'hash',
			type tinyint  NOT NULL DEFAULT '0' COMMENT '类型(1白名单,2黑名单)',
			dir varchar(8)  NOT NULL DEFAULT '' COMMENT '路径',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '备注',
			data_from tinyint NOT NULL DEFAULT '1' COMMENT '数据来源(1进程 2模块)',
			open_switch tinyint NOT NULL DEFAULT '0' COMMENT '开关',
            create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL,
			INDEX index_dir (dir) USING INVERTED,
			INDEX index_created_at (created_at) USING INVERTED
		)
		UNIQUE KEY(org_id,id) 
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

func (m *ProcGlobalModel) InsertItem(s *ProcGlobalModel) bool {

	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Create(s)
	})
	result := m.Exec(sql)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// GetItemByOrgAndId 通过企业id\hash获取一条记录
func (m *ProcGlobalModel) GetItemByOrgAndId(orgId int64, id string) (item *ProcGlobalModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND id=?", orgId, id).Take(&item)
	})
	result := m.Raw(m.TrimSql(sql)).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByOrgAndIds 通过企业id和终端id获取多条记录
func (m *ProcGlobalModel) GetItemsByOrgAndIds(orgId int64, ids []string) (items []*ProcGlobalModel, err error) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND id IN ?", orgId, ids).Find(&items)
	})
	result := m.Raw(m.TrimSql(sql)).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// DeleteByOrg 通过企业id删除
func (m *ProcGlobalModel) DeleteByOrg(orgId int64) bool {
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

// DeleteByOrgAndId 通过企业id、ID删除
func (m *ProcGlobalModel) DeleteByOrgAndId(orgId int64, id string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=? AND id=?", orgId, id).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}
