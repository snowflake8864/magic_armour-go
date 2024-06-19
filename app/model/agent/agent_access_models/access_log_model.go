package agent_access_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 访问日志表

func CreateAccessLogFactory() *AccessLogModel {
	return &AccessLogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type AccessLogModel struct {
	baseAccess
	model.BaseModel
	Id        int64  `json:"id" gorm:"primaryKey"`
	OrgId     int64  `json:"org_id"`
	AgentId   string `json:"agent_id"`
	Path      string `json:"path"`
	NetIp     string `json:"net_ip"`
	Input     string `json:"input"`
	Method    string `json:"method"`
	Status    int    `json:"status"`
	HttpCode  int    `json:"http_code"`
	DelayTime int64  `json:"delay_time"`
	Resp      string `json:"resp"`
	CreatedAt string `json:"created_at"`
}

// 表名
func (m *AccessLogModel) TableName() string {
	return "a_access_logs"
}

func (m *AccessLogModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			path varchar(128)  NOT NULL DEFAULT '' COMMENT '请求地址',
			created_at datetime NULL,
			id BIGINT NOT NULL DEFAULT '0',
			net_ip varchar(128)  NOT NULL DEFAULT '' COMMENT '通讯ip',
			input text NULL COMMENT '请求参数',
			method varchar(8) NULL COMMENT '请求方法',
			status tinyint NOT NULL DEFAULT '0' COMMENT '状态(1成功 0失败)',
			http_code smallint NOT NULL DEFAULT '0' COMMENT 'http状态码',
			delay_time int NOT NULL DEFAULT '0' COMMENT '延迟时间(ms)',
			resp string NOT NULL DEFAULT '' COMMENT '返回结果'
		)
		DUPLICATE  KEY(org_id, agent_id, path, created_at)
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

func (m *AccessLogModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *AccessLogModel) InsertItem(s *AccessLogModel) bool {

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

// DeleteByOrg 通过企业i删除
func (m *AccessLogModel) DeleteByOrg(orgId int64) bool {
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

// DeleteByOrgAndAgent 通过企业id和终端id删除
func (m *AccessLogModel) DeleteByOrgAndAgent(orgId int64, agentId string) bool {
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

// DeleteByOrgAndAgents 通通过企业id和终端id集合删除
func (m *AccessLogModel) DeleteByOrgAndAgents(orgId int64, agentIds []string) bool {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(m).Where("org_id=?", orgId).Where("agent_id IN ?", agentIds).Delete(m)
	})
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *AccessLogModel) DeleteByConditions(condition SearchConditions) bool {
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
