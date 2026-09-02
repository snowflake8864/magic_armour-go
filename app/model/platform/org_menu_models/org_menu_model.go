package org_menu_models

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// CreateOrgMenuFactory 企业用的菜单表
func CreateOrgMenuFactory() *OrgMenuModel {
	return &OrgMenuModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type OrgMenuModel struct {
	baseMenu
	model.BaseModel
	Id        int64  `gorm:"primaryKey" json:"id"`
	ParentId  int64  `json:"parent_id"`
	Position  int    `json:"position"`
	Type      int    `json:"type"`
	Name      string `json:"name"`
	EnName    string `json:"en_name"`
	Code      string `json:"code"`
	Sort      int    `json:"sort"`
	IsShow    int    `json:"is_show"`
	IsRoot    int    `json:"is_root"`
	UsageVer  int    `json:"usage_ver"`
	Remark    string `json:"remark"`
	Api       string `json:"api"`
	Level     int    `json:"level"`
	CreatedAt string `json:"created_at"`
}

// 表名

func (m *OrgMenuModel) TableName() string {
	return "p_org_menus"
}

func (m *OrgMenuModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			id int(11) unsigned NOT NULL AUTO_INCREMENT,
			parent_id int(11) NOT NULL DEFAULT '0' COMMENT '父级id',
			position tinyint(1) NOT NULL DEFAULT '1' COMMENT '方位(1左,2右,3顶部)',
			type tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型(1菜单,2按钮,3链接)',
			name varchar(128)  NOT NULL DEFAULT '' COMMENT '标题',
			code varchar(128)  NOT NULL DEFAULT '' COMMENT '页面标识',
			sort int(11) NOT NULL DEFAULT '1' COMMENT '排序',
			is_show tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否显示(0否,1是),前端判断用',
			is_root tinyint(1) NOT NULL DEFAULT '0' COMMENT '仅超管可用(0否,1是),后端判断用',
			usage_ver tinyint(1) NOT NULL DEFAULT '0' COMMENT '使用版本，二进制表示(0不区分版本,1离线版本,2在线版本)',
			remark varchar(512)  NOT NULL DEFAULT '' COMMENT '备注',
			api text  NOT NULL COMMENT '服务端对应接口(,)隔开',
			level smallint(4) NOT NULL DEFAULT '1' COMMENT '级别(从1开始)',
			created_at datetime NOT NULL,
			PRIMARY KEY (id),
			KEY index_parent (parent_id),
			KEY index_level (level)
		)ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) AlterTable(ver string) bool {
	sql := `ALTER TABLE ` + m.TableName()
	switch ver {
	case "MA_3.1.0_R4_B8":
		sql += `
			ADD COLUMN en_name varchar(256) NOT NULL DEFAULT '' COMMENT '英文名' AFTER name;
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

func (m *OrgMenuModel) TruncateTable() bool {
	sql := "TRUNCATE TABLE " + m.TableName()
	result := m.Exec(m.TrimSql(sql))
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) InitData() bool {
	_ = m.TruncateTable()

	/*
		(76, 0, 1, 1, '策略跳变', 'strategyJump_01', 187, 0, 0, 0, '', '', 1, '2025-03-12 00:00:00'),
		(77, 76, 1, 1, 'IP跳变', 'strategyJump_01_01', 1, 1, 0, 0, '', '', 2, '2025-03-12 00:00:00'),
		(78, 76, 1, 1, '口令跳变', 'strategyJump_01_02', 2, 1, 0, 0, '', '', 2, '2025-03-12 00:00:00'),
		(79, 76, 1, 1, '全局跳变配置', 'strategyJump_01_03', 3, 1, 0, 0, '', '', 2, '2025-03-12 00:00:00'),
		(80, 10, 1, 1, '备份回滚', 'agent_01_03', 3, 0, 0, 0, '', '', 2, '2025-03-12 00:00:00'),
	*/

	sql := `INSERT INTO ` + m.TableName()
	sql += ` VALUES `
	sql += `
		(1, 0, 1, 1, '首页', 'home_01', 100, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(2, 1, 1, 1, '概览', 'home_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(3, 0, 1, 1, '账户管理', 'account_01', 230, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(4, 3, 1, 1, '账户信息', 'account_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(5, 3, 1, 1, '账户安全', 'account_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(6, 0, 1, 1, '企业管理', 'business_01', 240, 1, 1, 0, '', '', 1, '2023-07-01 00:00:00'),
		(7, 6, 1, 1, '组织架构', 'business_01_01', 1, 1, 1, 0, '', '', 2, '2023-07-01 00:00:00'),
		(8, 0, 1, 1, '安全报告', 'safe_01', 260, 1, 1, 2, '云版', '', 1, '2023-07-01 00:00:00'),
		(9, 8, 1, 1, '安全通知', 'safe_01_01', 1, 1, 1, 2, '云版', '', 2, '2023-07-01 00:00:00'),
		(10, 0, 1, 1, '终端管理', 'agent_01', 110, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(11, 10, 1, 1, '终端管理', 'agent_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(12, 10, 1, 1, '全局信任目录', 'agent_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(13, 0, 1, 1, '可信防护', 'process_01', 120, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(14, 13, 1, 1, '进程管理', 'process_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(15, 13, 1, 1, '模块管理', 'process_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(17, 13, 1, 1, '进程全局名单', 'process_01_04', 4, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(18, 13, 1, 1, '模块全局名单', 'process_01_05', 5, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(19, 0, 1, 1, '勒索防护', 'black_01', 130, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(20, 19, 1, 1, '勒索防护', 'black_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(21, 19, 1, 1, '文件类型组设置', 'black_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(22, 0, 1, 1, '防篡改', 'tamper_01', 140, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(23, 22, 1, 1, '防篡改', 'tamper_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(25, 0, 1, 1, '端口虚开', 'port_01', 150, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(26, 25, 1, 1, '端口虚开', 'port_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(27, 25, 1, 1, '全局端口虚开', 'port_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(28, 0, 1, 1, '动态封堵', 'plug_01', 160, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(29, 28, 1, 1, '动态封堵', 'plug_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(30, 28, 1, 1, '全局封堵配置', 'plug_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(31, 28, 1, 1, '封堵全局名单', 'plug_01_03', 3, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(32, 28, 1, 1, '封堵列表', 'plug_01_04', 4, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(33, 0, 1, 1, '外设管理', 'usb_01', 170, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(34, 33, 1, 1, '外设管理', 'usb_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(35, 33, 1, 1, '外设全局名单', 'usb_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(36, 0, 1, 1, '外联管控', 'out_01', 180, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(37, 36, 1, 1, '外联管控', 'out_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(38, 36, 1, 1, '外联白名单', 'out_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(39, 0, 1, 1, '告警管理', 'alarm_01', 190, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(40, 39, 1, 1, '告警管理', 'alarm_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(41, 39, 1, 1, '进程异常告警', 'alarm_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(42, 39, 1, 1, '模块异常告警', 'alarm_01_03', 3, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(43, 39, 1, 1, '勒索告警', 'alarm_01_04', 4, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(44, 39, 1, 1, '防篡改告警', 'alarm_01_05', 5, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(45, 39, 1, 1, '端口虚开告警', 'alarm_01_06', 6, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(46, 39, 1, 1, '动态封堵告警', 'alarm_01_07', 7, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(47, 39, 1, 1, '外接设备告警', 'alarm_01_08', 8, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(48, 39, 1, 1, '外联管控告警', 'alarm_01_09', 9, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(49, 0, 1, 1, '授权管理', 'author_01', 210, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(50, 49, 1, 1, '终端授权', 'author_01_02', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(51, 0, 1, 1, '客户端管理', 'client_01', 220, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(52, 51, 1, 1, '下载页面', 'client_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(53, 0, 1, 1, '日志管理', 'log_01', 200, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(54, 53, 1, 1, '操作日志', 'log_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(55, 53, 1, 1, 'DNS日志', 'log_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(56, 53, 1, 1, '对外通讯日志', 'log_01_03', 3, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(57, 53, 1, 1, '对内通讯日志', 'log_01_04', 4, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(58, 53, 1, 1, '进程启动日志', 'log_01_05', 5, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(59, 53, 1, 1, '终端请求日志', 'log_01_06', 10, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(60, 53, 1, 1, '登录日志', 'log_01_07', 6, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(61, 39, 1, 1, '口令破解告警', 'alarm_01_10', 10, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(62, 0, 1, 1, '口令破解', 'ls_01', 185, 1, 0, 0, '', '', 1, '2023-07-01 00:00:00'),
		(63, 62, 1, 1, '口令破解', 'ls_01_01', 1, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(64, 62, 1, 1, '全局口令破解', 'ls_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(65, 49, 1, 1, '已完成订单', 'author_01_01', 1, 1, 0, 2, '云版', '', 2, '2023-07-01 00:00:00'),
		(66, 49, 1, 1, '购买记录', 'author_01_03', 3, 1, 1, 2, '云版', '', 2, '2023-07-01 00:00:00'),
		(67, 39, 1, 1, '自保告警', 'alarm_01_12', 20, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(68, 22, 1, 1, '全局防篡改', 'tamper_01_02', 2, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(69, 13, 1, 1, '提取样本', 'process_01_03', 3, 1, 1, 1, '离线版', '', 2, '2023-07-01 00:00:00'),
		(70, 6, 1, 1, '子用户管理', 'business_01_02', 2, 1, 1, 0, '', '', 2, '2023-07-01 00:00:00'),
		(71, 19, 1, 1, '全局勒索防护', 'black_01_03', 3, 1, 0, 0, '', '', 2, '2023-07-01 00:00:00'),
		(73, 36, 1, 1, '外联探测名单', 'out_01_03', 3, 1, 0, 0, '', '', 2, '2024-09-04 00:00:00'),
		(74, 49, 1, 1, '授权许可', 'author_01_04', 4, 1, 1, 1, '离线版', '', 2, '2024-12-16 00:00:00'),
		(75, 53, 1, 1, '终端上线日志', 'log_01_08', 6, 1, 0, 0, '', '', 2, '2025-01-14 00:00:00'),
		(81, 6, 1, 1, '权限管理', 'business_01_03', 3, 1, 1, 0, '', '', 2, '2025-03-26 00:00:00'),
		(82, 6, 1, 1, '角色管理', 'business_01_04', 4, 1, 1, 0, '', '', 2, '2025-03-26  00:00:00'),
		(83, 6, 1, 1, '安全管理', 'business_01_05', 5, 1, 1, 1, '', '', 2, '2025-07-23 00:00:00'),
		
		(84, 39, 1, 1, '病毒查杀告警', 'alarm_01_11', 11, 1, 0, 0, '', '', 2, '2026-02-26 00:00:00'),
		(85, 53, 1, 1, '漏洞扫描日志', 'log_01_09', 8, 1, 0, 0, '', '', 2, '2026-02-26 00:00:00');
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) InsertItem(s *OrgMenuModel) bool {
	result := m.Create(s)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) UpdateById(id int64, updates map[string]interface{}) bool {
	result := m.Table(m.TableName()).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) DeleteById(id int64) bool {
	result := m.Where("id=?", id).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

func (m *OrgMenuModel) GetItemById(id int64) (item *OrgMenuModel, err error) {
	result := m.Where("id=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return nil, result.Error
	}
}

func (m *OrgMenuModel) GetItemByParent(id int64) (item *OrgMenuModel, err error) {
	result := m.Where("parent_id=?", id).Take(&item)
	if result.Error == nil {
		return item, nil
	} else {
		return nil, result.Error
	}
}

// GetListByConditions 获取符合条件的全部记录
func (m *OrgMenuModel) GetListByConditions(condition SearchConditions) (items []*OrgMenuModel) {
	return m.GetItemsByConditions(condition, 0, 0, "")
}

// GetItemsByConditions 通过条件获取多条记录
func (m *OrgMenuModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order string) (items []*OrgMenuModel) {
	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	if limit > 0 {
		tmp = tmp.Offset(offset).Limit(limit)
	}
	if !string_helper.IsEmpty(order) {
		tmp = tmp.Order(order)
	} else {
		tmp = tmp.Order("level ASC, sort ASC")
	}
	result := tmp.Find(&items)
	if result.Error == nil {
		return items
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return nil
	}
}
