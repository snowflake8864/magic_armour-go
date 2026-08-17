package antivirus_engine_task_models

import (
	"gorm.io/gorm"
	"goskeleton/app/global/enums/enum_antivirus_engine_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 杀毒引擎升级任务表

func CreateUpgradeTaskFactory() *UpgradeTaskModel {
	return &UpgradeTaskModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type UpgradeTaskModel struct {
	baseOnline
	model.BaseModel
	Id                  int64   `json:"id" gorm:"primaryKey"`
	OrgId               int64   `json:"org_id"`
	Title               string  `json:"title"`
	TargetType          int     `json:"target_type"`
	TargetIds           string  `json:"target_ids"`
	EngineWindowsId     int64   `json:"engine_windows_id"`
	EngineLinuxId       int64   `json:"engine_linux_id"`
	UpgradeMode         int     `json:"upgrade_mode"`
	TaskType            int     `json:"task_type"`
	ExecStatus          int     `json:"exec_status"`
	ExecDate            string  `json:"exec_date"`
	ExecWeekly          int     `json:"exec_weekly"`
	ExecDaily           int     `json:"exec_daily"`
	ExecHis             string  `json:"exec_his"`
	ExecNextTs          int64   `json:"exec_next_ts"`
	ExecBatch           int     `json:"exec_batch"`
	ExecIsComplete      int     `json:"exec_is_complete"`
	ExecCompletedSource int     `json:"exec_completed_source"`
	ExecProgress        float64 `json:"exec_progress"`
	ExecAgentTotal      int     `json:"exec_agent_total"`
	ExecGoodTotal       int     `json:"exec_good_total"`
	ExecFailTotal       int     `json:"exec_fail_total"`
	ExecUnneededTotal   int     `json:"exec_unneeded_total"`
	ExecOfflineTotal    int     `json:"exec_offline_total"`
	ExecLiveTs          int64   `json:"exec_live_ts"`
	ExecStartAt         string  `json:"exec_start_at" gorm:"default:NULL"`
	ExecHandleAt        string  `json:"exec_handle_at" gorm:"default:NULL"`
	ExecCompleteAt      string  `json:"exec_complete_at" gorm:"default:NULL"`
	Remark              string  `json:"remark"`
	CreateId            int64   `json:"create_id"`
	CreateName          string  `json:"create_name"`
	CreatedAt           string  `json:"created_at"`
	UpdateId            int64   `json:"update_id"`
	UpdateName          string  `json:"update_name"`
	UpdatedAt           string  `json:"updated_at"`
}

// 表名
func (m *UpgradeTaskModel) TableName() string {
	return "a_antivirus_engine_upgrade_tasks"
}

func (m *UpgradeTaskModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id BIGINT(16) unsigned NOT NULL AUTO_INCREMENT,
			org_id int(11) NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			title varchar(128)  NOT NULL DEFAULT '' COMMENT '名称',
			target_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '目标类型',
			target_ids varchar(512) NOT NULL DEFAULT '0' COMMENT '目标ID集合',
			engine_windows_id int(11) NOT NULL DEFAULT '0' COMMENT '引擎：windows',
			engine_linux_id int(11) NOT NULL DEFAULT '0'  COMMENT '引擎：linux',
			upgrade_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '升级模式',
			task_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '任务类型',
			exec_status tinyint(1) NOT NULL DEFAULT '0' COMMENT '执行状态',
			exec_date date DEFAULT NULL COMMENT '指定日期执行',
			exec_weekly smallint(4) NOT NULL DEFAULT '0' COMMENT '周几执行',
			exec_daily smallint(4) NOT NULL DEFAULT '0' COMMENT '几号执行',
			exec_his varchar(16)  NOT NULL DEFAULT '' COMMENT '执行时分秒',
			exec_next_ts bigint(16) NOT NULL DEFAULT '0' COMMENT '下一次执行时间戳',
			exec_batch int(11) unsigned NOT NULL DEFAULT '0' COMMENT '执行批次',
			exec_is_complete smallint(4) NOT NULL DEFAULT '0' COMMENT '执行是否完成',
			exec_progress decimal(5, 2) NOT NULL DEFAULT '0.00' COMMENT '执行进度',
			exec_agent_total int(11) NOT NULL DEFAULT '0' COMMENT '执行终端总数',
			exec_good_total int(11) NOT NULL DEFAULT '0' COMMENT '执行成功终端总数',
			exec_fail_total int(11) NOT NULL DEFAULT '0' COMMENT '执行失败终端总数',
			exec_unneeded_total int(11) NOT NULL DEFAULT '0' COMMENT '执行无需终端总数',
			exec_offline_total int(11) NOT NULL DEFAULT '0' COMMENT '执行离线终端总数',
			exec_live_ts bigint(16) unsigned NOT NULL DEFAULT '0' COMMENT '执行活跃时间戳',
			exec_start_at datetime DEFAULT NULL COMMENT '执行开始时间',
			exec_handle_at datetime DEFAULT NULL COMMENT '执行处理时间',
			exec_complete_at datetime DEFAULT NULL COMMENT '执行结束时间',
			exec_completed_source smallint(4) NOT NULL DEFAULT '0' COMMENT '执行完成来源',
			remark text DEFAULT NULL  COMMENT '备注',
			create_id int(11) NOT NULL DEFAULT '0' COMMENT '创建人(user.id)',
			create_name varchar(128)  NOT NULL DEFAULT '' COMMENT '创建人(user.username)',
			created_at datetime NOT NULL,
			update_id int(11) NOT NULL DEFAULT '0' COMMENT '上一次编辑用户(user.id)',
			update_name varchar(128)  NOT NULL DEFAULT '' COMMENT '上一次编辑用户(user.username)',
			updated_at datetime NOT NULL,PRIMARY KEY (id),
			KEY index_org (org_id),
			KEY index_title (title),
			KEY index_target (target_type),
			KEY index_exec_status (exec_status,exec_next_ts,exec_start_at),
			KEY index_created_at (created_at)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='杀毒引擎:升级任务表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *UpgradeTaskModel) AlterTable(ver string) bool {
	sql := `ALTER TABLE ` + m.TableName()
	switch ver {
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

func (m *UpgradeTaskModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *UpgradeTaskModel) InsertItem(s *UpgradeTaskModel) bool {

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
func (m *UpgradeTaskModel) InsertBatch(list *[]UpgradeTaskModel) bool {
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

func (m *UpgradeTaskModel) conditions2OR(tx *gorm.DB, condition SearchConditions) {
	if !string_helper.IsEmpty(condition.TargetAgentId) && condition.TargetDeptId > 0 {
		tx.Where("(target_type=?) OR (target_type=? AND target_ids=?) OR (target_type=? AND FIND_IN_SET(?, target_ids))",
			enum_antivirus_engine_task.TargetTypeAllAgent,
			enum_antivirus_engine_task.TargetTypeSingleAgent, condition.TargetAgentId,
			enum_antivirus_engine_task.TargetTypeDepartments, condition.TargetDeptId)
	} else if !string_helper.IsEmpty(condition.TargetAgentId) {
		tx.Where("target_type=? OR (target_type=? AND target_ids1=?)",
			enum_antivirus_engine_task.TargetTypeAllAgent,
			enum_antivirus_engine_task.TargetTypeSingleAgent, condition.TargetAgentId)
	} else if condition.TargetDeptId > 0 {
		tx.Or("(target_type=? AND FIND_IN_SET(?, target_ids)",
			enum_antivirus_engine_task.TargetTypeAllAgent,
			enum_antivirus_engine_task.TargetTypeSingleAgent, condition.TargetDeptId)
	}
}

// DeleteByOrg 通过企业i删除
func (m *UpgradeTaskModel) DeleteByOrg(orgId int64) bool {
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

// DeleteByConditions 通过条件删除
func (m *UpgradeTaskModel) DeleteByConditions(condition SearchConditions) bool {
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

func (m *UpgradeTaskModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) bool {
	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemById 通过id获取一条记录
func (m *UpgradeTaskModel) GetItemById(id int64) (item *UpgradeTaskModel) {
	result := m.Model(m).Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item
	} else {
		return nil
	}
}

// GetListByConditions 通过条件获取多条记录
func (m *UpgradeTaskModel) GetListByConditions(condition SearchConditions) (items []*UpgradeTaskModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *UpgradeTaskModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*UpgradeTaskModel) {
	sql := m.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tmp := tx.Model(m)
		if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
			tmp = tmp.Where(whereSQL, vals...)
		}
		//OR条件
		m.conditions2OR(tmp, condition)

		if limit > 0 {
			tmp = tmp.Offset(offset).Limit(limit)
		}
		if !string_helper.IsEmpty(order) {
			tmp = tmp.Order(order)
		} else {
			tmp = tmp.Order("id DESC")
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
func (m *UpgradeTaskModel) Count(condition SearchConditions) int64 {
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
func (m *UpgradeTaskModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*UpgradeTaskModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}
