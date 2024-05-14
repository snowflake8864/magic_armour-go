package organization_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 企业日志构表

func CreateOrgLogFactory() *OrgLogModel {
	return &OrgLogModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OrgLogModel struct {
	model.BaseModel
	Id            int64  `json:"id" gorm:"primaryKey"`
	Module        int    `json:"module"`
	NetIp         string `json:"net_ip"`
	OrgId         int64  `json:"org_id"`
	OrgName       string `json:"org_name"`
	UserName      string `json:"user_name"`
	UserEmail     string `json:"user_email"`
	UserPhone     string `json:"user_phone"`
	Content       string `json:"content"`
	OperateSource int    `json:"operate_source"`
	OperateId     int64  `json:"operate_id"`
	OperateName   string `json:"operate_name"`
	CreatedAt     string `json:"created_at"`
}

// 表名
func (m *OrgLogModel) TableName() string {
	return "o_org_logs"
}

func (m *OrgLogModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
         id int(11) unsigned NOT NULL AUTO_INCREMENT,
		 module int(11) NOT NULL DEFAULT '0' COMMENT '业务模块',
         net_ip varchar(128) NOT NULL DEFAULT '' COMMENT '通讯ip',
		 org_id int(11) NOT NULL DEFAULT '0' COMMENT '企业(o_organization.id)',
		 org_name varchar(128) NOT NULL DEFAULT '' COMMENT '企业(o_organization.id)',
		 user_name varchar(128) NOT NULL DEFAULT '' COMMENT '默认管理用户(o_user.name)',
		 user_email varchar(64) NOT NULL DEFAULT '' COMMENT '默认管理用户(o_user.email)',
		 user_phone varchar(16) NOT NULL DEFAULT '' COMMENT '默认管理用户(o_user.phone)',
		 content text COMMENT '内容',
		 operate_source tinyint (1) NOT NULL DEFAULT '0' COMMENT '操作:来源',
		 operate_id int (11) NOT NULL DEFAULT '0' COMMENT '操作用户:编号',
         operate_name varchar (128) NOT NULL DEFAULT '' COMMENT '操作用户:名称',
         created_at datetime DEFAULT NULL,
		 PRIMARY KEY (id),
		 KEY module (module),
		 KEY org_id (org_id),
		 KEY org_name (org_name),
		 KEY user_name (user_name),
		 KEY user_phone (user_phone)
	   ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='企业:企业日志';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *OrgLogModel) InsertItem(s *OrgLogModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemsByConditions 通过条件获取多条记录
func (m *OrgLogModel) GetItemsByConditions(condition map[string]interface{}, offset, limit int, order string) (items []*OrgLogModel, err error) {
	tmp := m.Model(m)
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
	result := tmp.Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}

// Count 通过条件统计
func (m *OrgLogModel) Count(condition map[string]interface{}) int64 {
	var total int64
	tmp := m.Model(m).Select("id")
	if whereSQL, vals, _ := m.WhereBuild(condition); !helpers.StringIsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Count(&total)
	if result.Error == nil {
		return total
	} else {
		return 0
	}
}

// Page 通过条件分页
func (m *OrgLogModel) Page(condition map[string]interface{}, offset, limit int, order string) (total int64, items []*OrgLogModel) {
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
