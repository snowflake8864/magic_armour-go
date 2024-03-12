package model_platform

import (
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_login_solve"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers"
)

// 创建 数据字典表

func CreateDictionaryFactory() *DictionaryModel {
	return &DictionaryModel{BaseModel: model.BaseModel{DB: model.UseDbConn(variable.DbNamePlatform)}}
}

//结构体

type DictionaryModel struct {
	model.BaseModel
	Id     string `json:"id" gorm:"primaryKey"`
	Module string `json:"module"`
	Val    string `json:"val"`
	Remark string `json:"remark"`
}

// 表名

func (m *DictionaryModel) TableName() string {
	return "p_dictionary"
}

func (m *DictionaryModel) CreateTable() bool {
	sql := `CREATE TABLE IF NOT EXISTS ` + m.TableName()
	sql += `
		(
			Id varchar(64) NOT NULL DEFAULT '' COMMENT '标识',
			module varchar(32)  NOT NULL DEFAULT '' COMMENT '模块',
			val varchar(1024)  NOT NULL DEFAULT '' COMMENT '值',
			remark varchar(128)   NOT NULL DEFAULT '' COMMENT '备注',
			PRIMARY KEY (id),
			KEY index_module (module)
		) ENGINE=INNODB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据字典表';
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *DictionaryModel) InitData() bool {

	plugPortConfig := enum_plug.PlugPortConfig{}
	plugPortConfig.Format()
	jsonPlugPort := helpers.JsonEncode(plugPortConfig)

	plugLoginConfig := enum_plug.PlugLoginConfig{}
	plugLoginConfig.Format()
	jsonPlugLogin := helpers.JsonEncode(plugLoginConfig)

	loginSolveConfig := enum_login_solve.LoginSolveConfig{}
	loginSolveConfig.Format()
	jsonLoginSolve := helpers.JsonEncode(loginSolveConfig)

	//
	sql := `INSERT IGNORE INTO ` + m.TableName() + `(module, id, val, remark) VALUES`
	sql += `
			('admin', '` + enum_dictionary.AdminLoginExpire + `', '3','管理员:登录有效期(时)'),
			('admin', '` + enum_dictionary.AdminPwLengthMin + `', '6','管理员:密码长度最小值(6)'),
			('admin', '` + enum_dictionary.AdminPwLengthMax + `', '18','管理员:密码长度最小值(18)'),
			('admin', '` + enum_dictionary.AdminLoginErrorMax + `', '10','管理员:登录错误次数最大值(100)'),
			('admin', '` + enum_dictionary.AdminLoginErrorRemind + `', '3','管理员:登录错误次数提醒'),
			('admin', '` + enum_dictionary.AdminLoginLockDuration + `', '30','管理员:登录锁定时长(分)'),

			('agent', '` + enum_dictionary.AgentLogSendProto + `', '2','终端:日志发送协议(1udp 2tcp)'),
			('agent', '` + enum_dictionary.AgentSaveAccessLog + `', '0','终端：是否保存访问日志(0|1)'),
			('agent', '` + enum_dictionary.AgentGlobalPlugPortConfig + `', '` + jsonPlugPort + `','终端:全局封堵-端口虚开配置'),
			('agent', '` + enum_dictionary.AgentGlobalPlugLoginConfig + `', '` + jsonPlugLogin + `','终端：全局封堵-登录破解配置'),
			('agent', '` + enum_dictionary.AgentGlobalLoginSolveConfig + `', '` + jsonLoginSolve + `','终端：全局登录口令暴力破解告警配置'),

			('threaten', '` + enum_dictionary.ThreatenSwitch + `', '0','威胁情报:开关'),
			('threaten', '` + enum_dictionary.ThreatenUrl + `', '','威胁情报:服务器地址'),
			('threaten', '` + enum_dictionary.ThreatenAppId + `', '','威胁情报:应用id'),
			('threaten', '` + enum_dictionary.ThreatenAppKey + `', '','威胁情报:应用key'),
			('threaten', '` + enum_dictionary.ThreatenCacheExpire + `', '30','威胁情报:缓存结果保留周期(天)'),

			('huanshi', '` + enum_dictionary.HuanshiSwitch + `', '0','幻势:开关'),
			('huanshi', '` + enum_dictionary.HuanshiUrl + `', '','幻势:服务器地址'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitch + `', '0','幻势:告警开关'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmIp + `', '','幻势:告警ip'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmPort + `', '','幻势:告警端口'),
			('huanshi', '` + enum_dictionary.HuanshiLogSwitch + `', '0','幻势:日志开关'),
			('huanshi', '` + enum_dictionary.HuanshiLogIp + `', '','幻势:日志ip'),
			('huanshi', '` + enum_dictionary.HuanshiLogPort + `', '','幻势:日志端口'),
			('huanshi', '` + enum_dictionary.HuanshiDefaultVisitor + `', '2','幻势:默认访客'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitchProcess + `', '1','幻势:告警开关-进程'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitchDll + `', '1','幻势:告警开关-模块'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitchTamper + `', '1','幻势:告警开关-防篡改'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitchBlackmail + `', '1','幻势:告警开关-勒索'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitchPort + `', '1','幻势:告警开关-端口虚开'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitchUsb + `', '1','幻势:告警开关-外设'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitchPlug + `', '1','幻势:告警开关-动态封堵'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitchLogin + `', '1','幻势:告警开关-登录口令破解'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitchOutreach + `', '1','幻势:告警开关-外联'),
			('huanshi', '` + enum_dictionary.HuanshiAlarmSwitchSelf + `', '1','幻势:告警开关-自保'),

			('system', '` + enum_dictionary.SystemDeviceId + `', '','系统:设备标识码'),
			('system', '` + enum_dictionary.SystemPlatformProto + `', 'https','系统:平台访问协议'),
			('system', '` + enum_dictionary.SystemPlatformIp + `', '','系统:平台访问IP'),
			('system', '` + enum_dictionary.SystemPlatformPort + `', '443','系统:平台访问端口'),
			('system', '` + enum_dictionary.SystemAgentProto + `', 'https','系统:终端访问协议'),
			('system', '` + enum_dictionary.SystemAgentIp + `', '','系统:终端访问IP'),
			('system', '` + enum_dictionary.SystemAgentPort + `', '443','系统:终端访问端口'),
			('system', '` + enum_dictionary.SystemGrantAutoSwitch + `', '1','系统:自动授权开关'),
			('system', '` + enum_dictionary.SystemAlarmExpire + `', '7','系统:告警失效时间(天)，最大180'),
			('system', '` + enum_dictionary.SystemSyslogExpire + `', '7','系统:通信日志失效时间(天)，最大180'),
			('system', '` + enum_dictionary.SystemSshSwitch + `', '1','系统:ssh开关');
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *DictionaryModel) GetItemsByModule(module string) (items []DictionaryModel) {
	m.Where("module=?", module).Find(&items)
	return items
}

func (m *DictionaryModel) GetAll() (items []DictionaryModel) {
	m.Find(&items)
	return items
}

func (m *DictionaryModel) UpdateValById(id, val string) error {
	result := m.Model(&m).Where("id=?", id).Update("val", val)
	return result.Error
}
