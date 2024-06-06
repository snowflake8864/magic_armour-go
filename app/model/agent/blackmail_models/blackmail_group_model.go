package blackmail_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 表: 勒索防护文件类型组

func CreateBlackmailGroupFactory() *BlackmailGroupModel {
	return &BlackmailGroupModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type BlackmailGroupModel struct {
	model.BaseModel
	Id         int64  `json:"id" gorm:"primaryKey"`
	OrgId      int64  `json:"org_id"`
	Title      string `json:"title"`
	Suffix     string `json:"suffix"`
	CreateId   int64  `json:"create_id"`
	CreateName string `json:"create_name"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// 表名

func (m *BlackmailGroupModel) TableName() string {
	return "a_blackmail_group"
}

func (m *BlackmailGroupModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(	
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			id BIGINT NOT NULL DEFAULT '0',
			title varchar(128) NOT NULL DEFAULT '' COMMENT '标题',
			suffix varchar(128)  NOT NULL DEFAULT '' COMMENT '后缀,英文逗号隔开',
			create_id int NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			updated_at datetime NOT NULL
		)
		UNIQUE KEY(org_id, id) 
		DISTRIBUTED  BY HASH(org_id) BUCKETS AUTO 
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

func (m *BlackmailGroupModel) InsertItem(s *BlackmailGroupModel) bool {

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

// DeleteByOrg 通过企业id删除
func (m *BlackmailGroupModel) DeleteByOrg(orgId int64) bool {
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

// DeleteByOrgAndAndId 通过企业id、ID删除
func (m *BlackmailGroupModel) DeleteByOrgAndAndId(orgId int64, id string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Where("id=?", id).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByOrgAndAgentAndId 通过企业id、终端id、ID删除
func (m *BlackmailGroupModel) DeleteByOrgAndAgentAndId(orgId int64, agentId, id string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Where("agent_id=?", agentId).Where("id=?", id).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByOrgAndAgent 通过企业id和终端id删除
func (m *BlackmailGroupModel) DeleteByOrgAndAgent(orgId int64, agentId string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Where("agent_id=?", agentId).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}
