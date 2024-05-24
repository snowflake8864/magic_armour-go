package dir_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

//终端:目录结构

func CreateDirFactory() *DirModel {
	return &DirModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNameAgent)}}
}

//结构体

type DirModel struct {
	baseDir
	model.BaseModel
	OrgId      int64  `json:"org_id"`
	AgentId    string `json:"agent_id"`
	SearchDir  string `json:"search_dir"`
	FileName   string `json:"file_name"`
	FileType   int    `json:"file_type"`
	FileRw     string `json:"file_rw"`
	FileGroup  string `json:"file_group"`
	FileUser   string `json:"file_user"`
	FileSize   string `json:"file_size"`
	FileHash   string `json:"file_hash"`
	StartTime  string `json:"start_time"`
	UpdateTime string `json:"update_time"`
	CreatedAt  string `json:"created_at"`
}

// 表名

func (m *DirModel) TableName() string {
	return "a_dirs"
}

func (m *DirModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			org_id int NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			agent_id varchar(128) NOT NULL DEFAULT '' COMMENT '终端(agent.id)',
			search_dir varchar(512)  NOT NULL DEFAULT '' COMMENT '查询路径',
			created_at datetime NOT NULL,
			file_name varchar(512)  NOT NULL DEFAULT '' COMMENT '文件名',
			file_type tinyint NOT NULL DEFAULT '0' COMMENT '文件类型',
			file_rw  varchar(64)  NOT NULL DEFAULT '' COMMENT '权限',
			file_group  varchar(128)  NOT NULL DEFAULT '' COMMENT '用户组',
			file_user  varchar(128)  NOT NULL DEFAULT '' COMMENT '用户',
			file_size  varchar(128)  NOT NULL DEFAULT '' COMMENT '文件大小',
			file_hash  varchar(64)  NOT NULL DEFAULT '' COMMENT '文件hash',
			start_time  datetime NULL COMMENT '文件创建时间',
			update_time  datetime NULL COMMENT '文件修改时间'
	
		)
		DUPLICATE KEY(org_id,agent_id, search_dir, created_at)
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

// InsertBatch 插入多条
func (m *DirModel) InsertBatch(list *[]DirModel) bool {
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
func (m *DirModel) DeleteByConditions(condition SearchConditions) bool {
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
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return false
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *DirModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*DirModel) {
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
func (m *DirModel) Count(condition SearchConditions) int64 {
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
func (m *DirModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*DirModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
