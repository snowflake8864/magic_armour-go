package dictionary_models

import (
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/utils/helpers/json_helper"
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
			id varchar(64) NOT NULL DEFAULT '' COMMENT '标识',
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

	//
	sql := `INSERT IGNORE INTO ` + m.TableName() + `(module, id, val, remark) VALUES`
	sql += `
			('` + enum_dictionary.ModulePlatform + `', '` + enum_dictionary.PlatformLoginErrorMax + `', '10','管理员:登录错误次数最大值(100)'),
			('` + enum_dictionary.ModulePlatform + `', '` + enum_dictionary.PlatformLoginErrorRemind + `', '3','管理员:登录错误次数提醒'),
			('` + enum_dictionary.ModulePlatform + `', '` + enum_dictionary.PlatformLoginLockDuration + `', '30','管理员:登录锁定时长(分)'),
			('` + enum_dictionary.ModulePlatform + `', '` + enum_dictionary.PlatformSshSwitch + `', '1','管理员:系统ssh开关'),
			('` + enum_dictionary.ModulePlatform + `', '` + enum_dictionary.PlatformServiceCodeList + `', '["s0001"]','管理员:服务器编号列表'),

			('` + enum_dictionary.ModuleOrganization + `', '` + enum_dictionary.OrganizationLoginErrorMax + `', '10','企业:登录错误次数最大值(100)'),
			('` + enum_dictionary.ModuleOrganization + `', '` + enum_dictionary.OrganizationLoginErrorRemind + `', '3','企业:登录错误次数提醒'),
			('` + enum_dictionary.ModuleOrganization + `', '` + enum_dictionary.OrganizationLoginLockDuration + `', '30','企业:登录锁定时长(分)'),
			('` + enum_dictionary.ModuleOrganization + `', '` + enum_dictionary.OrganizationLoginSmsSwitch + `', '0','企业:登录短信验证开关'),
			('` + enum_dictionary.ModuleOrganization + `', '` + enum_dictionary.OrganizationSubUserLimit + `', '50','企业:子用户限制数目'),

			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentHttpUrl + `', 'https://192.168.4.1','终端:http请求地址(无端口)'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentHttpPort + `', '443','终端:http请求端口'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentLogSendProto + `', '2','终端:日志发送协议(1udp 2tcp)'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentSaveAccessLog + `', '1','终端：是否保存访问记录(0|1)'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentConfigExpire + `', '366','终端：配置文件有效期(天)'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentAutoDeleteExpire + `', '10','终端：自动删除有效期(天)'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentAlarmExpire + `', '10','终端：告警保留时长(天)'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentSyslogExpire + `', '10','终端：通讯日志保留时长(天)'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentProcGlobalSwitch + `', '1','终端：进程全局名单是否下发给终端开关'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentProcGlobalLimit + `', '2000','终端：进程全局名单限制数目'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentCompatibleDefaultOrg + `', '1','终端：兼容以前离线版本默认企业编号'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentCompatibleDefaultUser + `', '1','终端：兼容以前离线版本默认用户编号'),
			('` + enum_dictionary.ModuleAgent + `', '` + enum_dictionary.AgentTypes + `', '` + json_helper.Encode(enum_agent.AgentTypeDefault()) + `','终端：终端类型'),

			('` + enum_dictionary.ModuleThreaten + `', '` + enum_dictionary.ThreatenSwitch + `', '0','威胁情报:开关'),
			('` + enum_dictionary.ModuleThreaten + `', '` + enum_dictionary.ThreatenUrl + `', '','威胁情报:服务器地址'),
			('` + enum_dictionary.ModuleThreaten + `', '` + enum_dictionary.ThreatenAppId + `', '','威胁情报:应用id'),
			('` + enum_dictionary.ModuleThreaten + `', '` + enum_dictionary.ThreatenAppKey + `', '','威胁情报:应用key'),
			('` + enum_dictionary.ModuleThreaten + `', '` + enum_dictionary.ThreatenCacheExpire + `', '30','威胁情报:缓存结果保留周期(天)'),

			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiCallbackUrl + `', '0','幻势:回调地址(有端口)'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiCallIp + `', '0','幻势:回调ip'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiSwitch + `', '0','幻势:开关'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiUrl + `', '','幻势:服务器地址'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitch + `', '0','幻势:告警开关'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmIp + `', '','幻势:告警ip'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmPort + `', '514','幻势:告警端口'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiLogSwitch + `', '0','幻势:日志开关'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiLogIp + `', '','幻势:日志ip'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiLogPort + `', '514','幻势:日志端口'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitchProcess + `', '1','幻势:告警开关-进程'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitchDll + `', '1','幻势:告警开关-模块'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitchTamper + `', '1','幻势:告警开关-防篡改'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitchBlackmail + `', '1','幻势:告警开关-勒索'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitchPort + `', '1','幻势:告警开关-端口虚开'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitchUsb + `', '1','幻势:告警开关-外设'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitchPlug + `', '1','幻势:告警开关-动态封堵'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitchLogin + `', '1','幻势:告警开关-登录口令破解'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitchOutreach + `', '1','幻势:告警开关-外联'),
			('` + enum_dictionary.ModuleHuanshi + `', '` + enum_dictionary.HuanshiAlarmSwitchSelf + `', '1','幻势:告警开关-自保'),

			('` + enum_dictionary.ModuleOrder + `', '` + enum_dictionary.OrderBankcardNo + `', '3211 5010 0100 2884 12','订单： 银行卡-卡号'),
			('` + enum_dictionary.ModuleOrder + `', '` + enum_dictionary.OrderBankcardName + `', '北京卫达信息技术有限公司','订单： 银行卡-卡名'),
			('` + enum_dictionary.ModuleOrder + `', '` + enum_dictionary.OrderBankcardBank + `', '兴业银行北京月坛支行','订单： 银行卡-开户行'),
			('` + enum_dictionary.ModuleOrder + `', '` + enum_dictionary.OrderPayingTel + `', '010-53516777','订单： 支付时联系电话'),
			('` + enum_dictionary.ModuleOrder + `', '` + enum_dictionary.OrderExpireDays + `', '7','订单： 过期天数');

		`
	result := m.Exec(m.TrimSql(sql))
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func (m *DictionaryModel) GetItemsByModule(module string) (items []*DictionaryModel) {
	m.Where("module=?", module).Find(&items)
	return items
}

func (m *DictionaryModel) GetAll() (items []*DictionaryModel) {
	m.Find(&items)
	return items
}

func (m *DictionaryModel) UpdateValById(id string, val interface{}) error {
	result := m.Model(&m).Where("id=?", id).Update("val", val)
	return result.Error
}
