package organization_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
)

// 企业扩展表

func CreateExtendFactory() *ExtendModel {
	return &ExtendModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type ExtendModel struct {
	baseOrg
	model.BaseModel
	Id                   int64  `json:"id" gorm:"primaryKey"`
	PlugPortConfig       string `json:"plug_port_config"`
	PlugLoginConfig      string `json:"plug_login_config"`
	LoginSolveConfig     string `json:"login_solve_config"`
	AlarmNoticeConfig    string `json:"alarm_notice_config"`
	HardwareNoticeConfig string `json:"hardware_notice_config"`
	PwJumpConfig         string `json:"pw_jump_config"`
	IpJumpConfig         string `json:"ip_jump_config"`
	UpdatedAt            string `json:"updated_at"`
}

// 表名
func (m *ExtendModel) TableName() string {
	return "o_org_extend"
}

func (m *ExtendModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id INT(11) UNSIGNED NOT NULL COMMENT '同(a_organizations.id)',
			plug_port_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵端口虚开策略(json)',
			plug_login_config varchar(256)  NOT NULL DEFAULT '' COMMENT '封堵口令破解策略(json)',
			login_solve_config varchar(256)  NOT NULL DEFAULT '' COMMENT '口令破解策略(json)',
			alarm_notice_config VARCHAR(512) DEFAULT '' COMMENT '告警通知策略(json)',
			hardware_notice_config VARCHAR(256) DEFAULT '' COMMENT '硬件通知策略(json)',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
        ) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC COMMENT='企业:扩展表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *ExtendModel) AlterTable(ver string) bool {
	sql := `ALTER TABLE ` + m.TableName()
	switch ver {
	case "MAT_3.1.0_R1":
		sql += `
			ADD COLUMN pw_jump_config varchar(256) NOT NULL DEFAULT '' COMMENT '口令跳变:策略(json)' AFTER hardware_notice_config,
			ADD COLUMN ip_jump_config varchar(256) NOT NULL DEFAULT '' COMMENT 'ip跳变:策略(json)' AFTER pw_jump_config;
		`
		break
	default:
		return true
	}
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

// InsertItem 插入一条
func (m *ExtendModel) InsertItem(s *ExtendModel) bool {
	result := m.Model(m).Create(s)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// UpdateById 通过id修改
func (m *ExtendModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Model(m).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *ExtendModel) DeleteById(id int64) bool {
	result := m.Where("id=?", id).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *ExtendModel) GetItemById(id int64) (item *ExtendModel, err error) {
	result := m.Model(m).Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return item, result.Error
	}
}

// GetItemsByIds 通过id集合获取多条记录
func (m *ExtendModel) GetItemsByIds(ids []int64) (items []*ExtendModel, err error) {
	result := m.Model(m).Where("id IN ?", ids).Find(&items)
	if result.Error == nil {
		return items, nil
	} else {
		return nil, result.Error
	}
}
