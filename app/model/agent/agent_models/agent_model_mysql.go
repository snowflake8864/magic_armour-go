package agent_models

import (
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/string_helper"
)

// 终端表

func CreateAgentFactory() *AgentModel {
	return &AgentModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

// 结构体
type AgentModel struct {
	baseAgent
	model.BaseModel
	AgentId             string `json:"id"  gorm:"primaryKey;column:id"`
	OrgId               int64  `json:"org_id"`
	HostName            string `json:"host_name"`
	MacId               string `json:"mac_id"`
	AgentIp             string `json:"agent_ip"`
	NetIp               string `json:"net_ip"`
	OsName              string `json:"os_name"`
	PlatformType        int    `json:"platform_type"`
	OsType              int    `json:"os_type"`
	ArchType            int    `json:"arch_type"`
	MemDesc             string `json:"mem_desc"`
	HdDesc              string `json:"hd_desc"`
	CpuDesc             string `json:"cpu_desc"`
	NetworkCard         string `json:"network_card"`
	AgentStartAt        string `json:"agent_start_at"`
	OsStartAt           string `json:"os_start_at"`
	Ver                 string `json:"ver"`
	ModVer              string `json:"mod_ver"`
	Auth                string `json:"auth"`
	IsOnline            int    `json:"is_online"`
	OnlineAt            string `json:"online_at"`
	LiveAt              string `json:"live_at"`
	IsSafe              int    `json:"is_safe"`
	UpgradeVer          string `json:"upgrade_ver"`
	UserId              int64  `json:"user_id"`
	DeptId              int64  `json:"dept_id"`
	AgentType           int    `json:"agent_type"`
	Alias               string `json:"alias"`
	ContactName         string `json:"contact_name"`
	ContactTel          string `json:"contact_tel"`
	Position            string `json:"position"`
	Remark              string `json:"remark"`
	GrantConfig         string `json:"grant_config"`
	GrantOrderNo        string `json:"grant_order_no"`
	GrantProductTitle   string `json:"grant_product_title"`
	GrantExpireAt       string `json:"grant_expire_at"  gorm:"default:NULL"`
	GrantType           int    `json:"grant_type"`
	GrantStatus         int    `json:"grant_status"`
	ProcessSwitch       int    `json:"process_switch"`
	ProcessMode         int    `json:"process_mode"`
	DllSwitch           int    `json:"dll_switch"`
	DllMode             int    `json:"dll_mode"`
	TamperSwitch        int    `json:"tamper_switch"`
	TamperMode          int    `json:"tamper_mode"`
	BlackmailSwitch     int    `json:"blackmail_switch"`
	BlackmailMode       int    `json:"blackmail_mode"`
	PortSwitch          int    `json:"port_switch"`
	PlugSwitch          int    `json:"plug_switch"`
	UsbSwitch           int    `json:"usb_switch"`
	UsbMode             int    `json:"usb_mode"`
	SyslogDnsSwitch     int    `json:"syslog_dns_switch"`
	SyslogOuterSwitch   int    `json:"syslog_outer_switch"`
	SyslogInnerSwitch   int    `json:"syslog_inner_switch"`
	SyslogProcessSwitch int    `json:"syslog_process_switch"`
	SyslogLoginSwitch   int    `json:"syslog_login_switch"`
	OutreachAlarmSwitch int    `json:"outreach_alarm_switch"`
	SelfProtectSwitch   int    `json:"self_protect_switch"`
	KernelUpgradeSwitch int    `json:"kernel_upgrade_switch"`
	InternetSwitch      int    `json:"internet_switch"`
	CommonTime          int    `json:"common_time"`
	FastTime            int    `json:"fast_time"`
	AutoWhiteSwitch     int    `json:"auto_white_switch"`
	AlarmAiSwitch       int    `json:"alarm_ai_switch"`
	LoginSolveSwitch    int    `json:"login_solve_switch"`
	HardwareSwitch      int    `json:"hardware_switch"`
	HardwareTime        int    `json:"hardware_time"`
	PwJumpSwitch        int    `json:"pw_jump_switch"`
	IpJumpSwitch        int    `json:"ip_jump_switch"`
	CreatedAt           string `json:"created_at"`
}

// 表名

func (m *AgentModel) TableName() string {
	return "a_agents"
}

func (m *AgentModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
            org_id int(11) NOT NULL DEFAULT '0' COMMENT '企业(organization.id)',
			id varchar(128) NOT NULL DEFAULT '' COMMENT '设备id',
			host_name varchar(64)  NOT NULL DEFAULT '' COMMENT '主机名',
			mac_id varchar(512)  NOT NULL DEFAULT '' COMMENT '物理地址',
			agent_ip varchar(512)  NOT NULL DEFAULT '' COMMENT '终端ip',
			net_ip varchar(64)  NOT NULL DEFAULT '' COMMENT '通讯ip',
			os_name varchar(64)  NOT NULL DEFAULT '' COMMENT '系统名称',
			platform_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '平台类型(win|linux, 统一定义)',
			os_type smallint(4) NOT NULL DEFAULT '0' COMMENT '系统类别(参考ubuntu|centos等,统一定义)',
			mem_desc varchar(128)  NOT NULL DEFAULT '' COMMENT '内存',
			hd_desc varchar(128)  NOT NULL DEFAULT '' COMMENT '硬盘',
			cpu_desc varchar(128)  NOT NULL DEFAULT '' COMMENT '处理器',
			network_card varchar(256)  NOT NULL DEFAULT '' COMMENT '网卡',
			agent_start_at datetime NOT NULL COMMENT '终端启动时间',
			os_start_at datetime NOT NULL COMMENT '系统启动时间',
			ver varchar(32)  NOT NULL DEFAULT '' COMMENT '终端版本号',
			auth varchar(16)  NOT NULL DEFAULT '' COMMENT '安全密码，终端上传',
			is_online tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否在线(0不在,1在)',
			online_at datetime NOT NULL COMMENT '上线开始时间',
			live_at datetime NOT NULL COMMENT '最近活动时间',
			is_safe tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否安全(0不安全,1安全)',
			upgrade_ver varchar(32)  NOT NULL DEFAULT '' COMMENT '要升级版本号',
			user_id int(11) NOT NULL DEFAULT '0' COMMENT '管理员编号(user.id)',
			dept_id int(11) NOT NULL DEFAULT '0' COMMENT '组织架构编号(department.id)',
			agent_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '终端类型',
			alias varchar(64)  NOT NULL DEFAULT '' COMMENT '别名',
			contact_name varchar(32)  NOT NULL DEFAULT '' COMMENT '责任人:姓名',
			contact_tel varchar(32)  NOT NULL DEFAULT '' COMMENT '责任人:电话',
			position varchar(128)  NOT NULL DEFAULT '' COMMENT '位置',
			remark varchar(128)  NOT NULL DEFAULT '' COMMENT '备注',
			grant_config text NOT NULL COMMENT '授权:配置(json)',
			grant_order_no varchar(64)  NOT NULL DEFAULT '' COMMENT '授权:订单号',
			grant_product_title varchar(128)  NOT NULL DEFAULT '' COMMENT '授权:产品名称',
			grant_expire_at date NULL COMMENT '授权:过期时间',
			grant_type tinyint(1) NOT NULL DEFAULT '0' COMMENT '授权:类型',
			grant_status tinyint(1) NOT NULL DEFAULT '0' COMMENT '授权:状态',
			process_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '进程:开关',
			process_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '进程:模式',
			dll_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '模块:开关',
			dll_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '模块:模式',
			tamper_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '防篡改:开关',
			tamper_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '防篡改:模式',
			blackmail_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '勒索:开关',
			blackmail_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '勒索:模式',
			port_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '端口虚开:开关',
			plug_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '封堵:开关',
			usb_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '外设:开关',
			usb_mode tinyint(1) NOT NULL DEFAULT '0' COMMENT '外设:模式',
			syslog_dns_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '通信日志:dns开关',
			syslog_outer_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '通信日志:对外通讯开关',
			syslog_inner_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '通信日志:对内通讯开关',
			syslog_process_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '通信日志:进程启动开关',
			syslog_login_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '通信日志:登录日志开关',
			outreach_alarm_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '外联告警:开关',
			self_protect_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '自我保护:开关',
			kernel_upgrade_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '内核升级:开关',
			common_time int(11) NOT NULL DEFAULT '60' COMMENT ' 心跳间隔:普通(秒)',
			fast_time int(11) NOT NULL DEFAULT '5' COMMENT '心跳间隔:快速(秒)',
			auto_white_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '自动匹配白名单:开关',
			alarm_ai_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '告警智能处理:开关',
			login_solve_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '口令破解:开关',
			hardware_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '上传硬件资源:开关',
			hardware_time int(11) NOT NULL DEFAULT '300' COMMENT '上传硬件资源:快速(秒)',
			created_at datetime NOT NULL,
	    	PRIMARY KEY (org_id, id),
		    KEY index_agent_ip (agent_ip),
			KEY index_dept_id (dept_id) ,
			KEY index_created_at (created_at)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='终端表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *AgentModel) AlterTable(ver string) bool {
	sql := `ALTER TABLE ` + m.TableName()
	switch ver {
	case "3.1.0_T4":
		sql += `
			ADD COLUMN mod_ver varchar(64) NOT NULL DEFAULT '' COMMENT '终端内核版本' AFTER ver,
			ADD COLUMN internet_switch tinyint(1) NOT NULL DEFAULT '1' COMMENT '上网准入:开关' AFTER kernel_upgrade_switch;
		`
	case "MAT_3.1.0_R1":
		sql += `
			ADD COLUMN pw_jump_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT '密码跳变:开关' AFTER hardware_time,
			ADD COLUMN ip_jump_switch tinyint(1) NOT NULL DEFAULT '0' COMMENT 'ip跳变:开关' AFTER pw_jump_switch;
		`
	case "MA_3.1.0_R4_B5":
		sql += `
			ADD COLUMN arch_type tinyint(1) NOT NULL DEFAULT '1' COMMENT 'CPU架构(参考amd64|arm64等,统一定义)' AFTER os_type;
		`
	case "MA_3.1.0_R4_B8":
		sql += `
			MODIFY COLUMN os_name varchar(128) NOT NULL DEFAULT '' COMMENT '系统名称' AFTER net_ip;
		`
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
func (m *AgentModel) InsertItem(s *AgentModel) bool {
	result := m.Create(s)
	if result.RowsAffected > 0 {
		return true
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return false
	}
}

// InsertBatch 插入多条
func (m *AgentModel) InsertBatch(list *[]AgentModel) bool {
	result := m.Create(list)
	if result.RowsAffected > 0 {
		return true
	} else {
		variable.ZapLog.Sugar().Error(result.Error.Error())
		return false
	}
}

// UpdateById 通过id修改
func (m *AgentModel) UpdateById(id string, updates map[string]interface{}) (bool, int64) {
	result := m.Table(m.TableName()).Where("id=?", id).Updates(updates)
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// UpdateByIds 通过id集合修改
func (m *AgentModel) UpdateByIds(ids []string, updates map[string]interface{}) (bool, int64) {
	result := m.Table(m.TableName()).Where("id IN ?", ids).Updates(updates)
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// UpdateByOrgAndAgent 通过企业id和终端id修改
func (m *AgentModel) UpdateByOrgAndAgent(orgId int64, agentId string, updates map[string]interface{}) (bool, int64) {
	result := m.Model(m).Where("org_id=? AND id=?", orgId, agentId).Updates(updates)
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// UpdateByOrgAndAgents 通过企业id和终端id集合修改
func (m *AgentModel) UpdateByOrgAndAgents(orgId int64, agentIds []string, updates map[string]interface{}) (bool, int64) {
	result := m.Model(m).Where("org_id=? AND id IN ?", orgId, agentIds).Updates(updates)
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// UpdateByConditions 通过条件修改
func (m *AgentModel) UpdateByConditions(condition SearchConditions, updates map[string]interface{}) (bool, int64) {
	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Updates(updates)
	if result.Error == nil {
		return true, result.RowsAffected
	} else {
		return false, 0
	}
}

// DeleteByOrgAndAgent 通过企业id和终端id删除
func (m *AgentModel) DeleteByOrgAndAgent(orgId int64, agentId string) bool {
	result := m.Model(m).Where("org_id=? AND id=?", orgId, agentId).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByOrg 通过企业i删除
func (m *AgentModel) DeleteByOrg(orgId int64) bool {
	result := m.Model(m).Where("org_id=?", orgId).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByOrgAndAgents 通通过企业id和终端id集合删除
func (m *AgentModel) DeleteByOrgAndAgents(orgId int64, agentIds []string) bool {
	result := m.Model(m).Where("org_id=? AND id IN ?", orgId, agentIds).Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// DeleteByConditions 通过条件删除
func (m *AgentModel) DeleteByConditions(condition SearchConditions) bool {
	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Delete(m)
	if result.Error == nil {
		return true
	} else {
		return false
	}
}

// GetItemByOrgAndAgent 通过企业id和终端id获取一条记录
func (m *AgentModel) GetItemByOrgAndAgent(orgId int64, agentId, column string) (item *AgentModel) {
	result := m.Model(m).Select(m.GetColumn(column)).Where("org_id=? AND id=?", orgId, agentId).Take(&item)
	if result.Error == nil {
		return item
	} else {
		return item
	}
}

// GetItemsByOrg 通过企业id获取多条记录
func (m *AgentModel) GetItemsByOrg(orgId int64, column string) (items []*AgentModel) {
	result := m.Model(m).Select(m.GetColumn(column)).Where("org_id=?", orgId).Find(&items)
	if result.Error == nil {
		return items
	} else {
		return nil
	}
}

// GetMapsByOrgAndAgents 通过企业id和终端ids获取多条记录, 有映射关系
func (m *AgentModel) GetMapsByOrgAndAgents(orgId int64, agentIds []string, column string) map[string]*AgentModel {
	models := m.GetItemsByOrgAndAgents(orgId, agentIds, column)

	if models == nil {
		return nil
	}
	result := make(map[string]*AgentModel)
	for _, item := range models {
		result[item.AgentId] = item
	}
	return result
}

// GetItemsByOrgAndAgents 通过企业id和终端id获取多条记录
func (m *AgentModel) GetItemsByOrgAndAgents(orgId int64, agentIds []string, column string) (items []*AgentModel) {
	result := m.Model(m).Select(m.GetColumn(column)).Where("org_id=? AND id IN ?", orgId, agentIds).Find(&items)
	if result.Error == nil {
		return items
	} else {
		return nil
	}
}

// GetIdsByByConditions s 通过条件获取ID集合
func (m *AgentModel) GetIdsByByConditions(condition SearchConditions) (ids []string) {
	models := m.GetListByConditions(condition, m.ColumnBasic(true, false))
	if models == nil {
		return ids
	}

	for _, item := range models {
		ids = append(ids, item.AgentId)
	}
	return ids
}

// GetListByConditions 通过条件获取多条记录
func (m *AgentModel) GetListByConditions(condition SearchConditions, column string) (items []*AgentModel) {
	return m.GetItemsByConditions(condition, 0, 0, "", column)
}

// GetItemsByConditions 通过条件获取多条记录
func (m *AgentModel) GetItemsByConditions(condition SearchConditions, offset, limit int, order, column string) (items []*AgentModel) {

	tmp := m.Model(m)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	if limit > 0 {
		tmp = tmp.Offset(offset).Limit(limit)
	}
	if order == enum_agent.AgentFieldDefault {
		tmp = tmp.Order("is_online DESC").Order("is_safe DESC").Order("created_at ASC")
	} else if !string_helper.IsEmpty(order) {
		tmp = tmp.Order(order)
	} else {
		tmp = tmp.Order("created_at ASC")
	}
	result := tmp.Find(&items)
	if result.Error == nil {
		return items
	} else {
		return nil
	}

}

// Count 通过条件统计
func (m *AgentModel) Count(condition SearchConditions) int64 {
	var total int64
	tmp := m.Model(m).Select("id")
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
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
func (m *AgentModel) Page(condition SearchConditions, offset, limit int, order string) (total int64, items []*AgentModel) {
	if total = m.Count(condition); total == 0 {
		return total, nil
	}
	if items = m.GetItemsByConditions(condition, offset, limit, order, "*"); len(items) > 0 {
		return total, items
	} else {
		return 0, nil
	}
}

// CalcByConditions 通过条件获取多条记录
func (m *AgentModel) CalcByConditions(condition SearchConditions, column string) (item *AgentCalc) {
	tmp := m.Model(m).Select(column)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Take(&item)
	if result.Error == nil {
		return item
	} else {
		return nil
	}

}

// CalcByConditionsGroupOrg 通过条件获取多条记录
func (m *AgentModel) CalcByConditionsGroupOrg(condition SearchConditions, column string) (items []*AgentCalc) {
	tmp := m.Model(m).Select(column)
	if whereSQL, vals, _ := m.WhereBuild(m.conditions2Where(condition)); !string_helper.IsEmpty(whereSQL) {
		tmp = tmp.Where(whereSQL, vals...)
	}
	result := tmp.Group("org_id").Find(&items)
	if result.Error == nil {
		return items
	} else {
		return nil
	}
}

// CalcByConditionsGroupAgentReturnMap 通过条件获取多条记录且按企业分组
func (m *AgentModel) CalcByConditionsGroupAgentReturnMap(condition SearchConditions, column string) map[int64]*AgentCalc {
	var result = make(map[int64]*AgentCalc)

	items := m.CalcByConditionsGroupOrg(condition, column)
	if items == nil || len(items) == 0 {
		return result
	}
	for _, item := range items {
		result[item.OrgId] = item
	}

	return result
}
