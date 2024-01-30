package modelPlatform

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
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

	sql := `INSERT IGNORE INTO ` + m.TableName() + `(module, id, val, remark) VALUES`
	sql += `
			('admin', '` + enums.DictionaryAdminLoginExpire + `', '3','管理员:登录有效期(时)'),
			('admin', '` + enums.DictionaryAdminPwLengthMin + `', '6','管理员:密码长度最小值(6)'),
			('admin', '` + enums.DictionaryAdminPwLengthMax + `', '18','管理员:密码长度最小值(18)'),
			('admin', '` + enums.DictionaryAdminLoginErrorMax + `', '10','管理员:登录错误次数最大值(100)'),
			('admin', '` + enums.DictionaryAdminLoginErrorRemind + `', '3','管理员:登录错误次数提醒'),
			('admin', '` + enums.DictionaryAdminLoginLockDuration + `', '30','管理员:登录锁定时长(分)'),

			('agent', '` + enums.DictionaryAgentLogSendProto + `', '2','终端:日志发送协议(1udp 2tcp)'),
			('agent', '` + enums.DictionaryAgentSaveAccessLog + `', '0','终端：是否保存访问日志(0|1)'),
			('agent', '` + enums.DictionaryAgentGlobalPlugPortConfig + `', '','终端:全局封堵-端口虚开配置'),
			('agent', '` + enums.DictionaryAgentGlobalPlugLoginConfig + `', '','终端：全局封堵-登录破解配置'),
			('agent', '` + enums.DictionaryAgentGlobalLoginSolveConfig + `', '','终端：全局登录口令暴力破解告警配置'),

			('threaten', '` + enums.DictionaryThreatenSwitch + `', '0','威胁情报:开关'),
			('threaten', '` + enums.DictionaryThreatenUrl + `', '','威胁情报:服务器地址'),
			('threaten', '` + enums.DictionaryThreatenAppId + `', '','威胁情报:应用id'),
			('threaten', '` + enums.DictionaryThreatenAppKey + `', '','威胁情报:应用key'),
			('threaten', '` + enums.DictionaryThreatenCacheExpire + `', '30','威胁情报:缓存结果保留周期(天)'),

			('huanshi', '` + enums.DictionaryHuanshiSwitch + `', '0','幻势:开关'),
			('huanshi', '` + enums.DictionaryHuanshiUrl + `', '','幻势:服务器地址'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitch + `', '0','幻势:告警开关'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmIp + `', '','幻势:告警ip'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmPort + `', '','幻势:告警端口'),
			('huanshi', '` + enums.DictionaryHuanshiLogSwitch + `', '0','幻势:日志开关'),
			('huanshi', '` + enums.DictionaryHuanshiLogIp + `', '','幻势:日志ip'),
			('huanshi', '` + enums.DictionaryHuanshiLogPort + `', '','幻势:日志端口'),
			('huanshi', '` + enums.DictionaryHuanshiDefaultVisitor + `', '2','幻势:默认访客'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitchProcess + `', '1','幻势:告警开关-进程'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitchDll + `', '1','幻势:告警开关-模块'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitchTamper + `', '1','幻势:告警开关-防篡改'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitchBlackmail + `', '1','幻势:告警开关-勒索'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitchPort + `', '1','幻势:告警开关-端口虚开'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitchUsb + `', '1','幻势:告警开关-外设'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitchPlug + `', '1','幻势:告警开关-动态封堵'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitchLogin + `', '1','幻势:告警开关-登录口令破解'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitchOutreach + `', '1','幻势:告警开关-外联'),
			('huanshi', '` + enums.DictionaryHuanshiAlarmSwitchSelf + `', '1','幻势:告警开关-自保'),

			('system', '` + enums.DictionarySystemDeviceId + `', '1','系统:设备标识码'),
			('system', '` + enums.DictionarySystemPlatformUrl + `', '1','系统:平台访问地址'),
			('system', '` + enums.DictionarySystemPlatformIp + `', '1','系统:平台IP地址'),
			('system', '` + enums.DictionarySystemPlatformPort + `', '443','系统:平台开放端口'),
			('system', '` + enums.DictionarySystemAgentUrl + `', '','系统:终端使用的访问地址'),
			('system', '` + enums.DictionarySystemAgentIp + `', '','系统:终端使用的IP地址'),
			('system', '` + enums.DictionarySystemAgentPort + `', '443','系统:终端使用的开放端口'),
			('system', '` + enums.DictionarySystemGrantAutoSwitch + `', '1','系统:自动授权开关'),
			('system', '` + enums.DictionarySystemAlarmExpire + `', '1','系统:告警失效时间(天)，最大180'),
			('system', '` + enums.DictionarySystemSyslogExpire + `', '1','系统:通信日志失效时间(天)，最大180'),
			('system', '` + enums.DictionarySystemSshSwitch + `', '1','系统:ssh开关');
		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
