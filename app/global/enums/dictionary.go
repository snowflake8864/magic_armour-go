package enums

// 数据字典

const (
	DictionaryModuleAdmin    string = "admin"    //模块：管路员
	DictionaryModuleAgent    string = "agent"    //模块：终端
	DictionaryModuleThreaten string = "threaten" //模块：终端
	DictionaryModuleHuanshi  string = "huanshi"  //模块：幻势
	DictionaryModuleSystem   string = "system"   //模块：系统

	DictionaryAdminLoginExpire       string = "admin_login_expire"        //管理员:登录有效期(时)
	DictionaryAdminPwLengthMin       string = "admin_pw_length_min"       //管理员:密码长度最小值(6)
	DictionaryAdminPwLengthMax       string = "admin_pw_length_max"       //管理员:密码长度最小值(18)
	DictionaryAdminLoginErrorMax     string = "admin_login_error_max"     //管理员:登录错误次数最大值(100)
	DictionaryAdminLoginErrorRemind  string = "admin_login_error_remind"  //管理员:登录错误次数提醒
	DictionaryAdminLoginLockDuration string = "admin_login_lock_duration" // 管理员:登录锁定时长(分)

	DictionaryAgentLogSendProto           string = "agent_log_send_proto"            //终端:日志发送协议(1udp 2tcp)
	DictionaryAgentSaveAccessLog          string = "agent_save_access_log"           //终端：是否保存访问日志(0|1)
	DictionaryAgentTaskCacheExpire        string = "agent_task_cache_expire"         //终端:任务缓存过期时间(时)
	DictionaryAgentGlobalPlugConfig       string = "agent_global_plug_port_config"   //终端:全局封堵-端口虚开配置
	DictionaryAgentGlobalPlugLoginConfig  string = "agent_global_plug_login_config"  //终端：全局封堵-登录破解配置
	DictionaryAgentGlobalLoginSolveConfig string = "agent_global_Login_solve_config" //终端：全局登录口令暴力破解告警配置

	DictionaryThreatenSwitch      string = "threaten_switch"       //威胁情报:开关
	DictionaryThreatenUrl         string = "threaten_url"          //威胁情报:服务器地址
	DictionaryThreatenAppId       string = "threaten_app_id"       //威胁情报:应用id
	DictionaryThreatenAppKey      string = "threaten_app_key"      //威胁情报:应用key
	DictionaryThreatenCacheExpire string = "threaten_cache_expire" //威胁情报:缓存结果保留周期(天)

	DictionaryHuanshiTokenExpire          string = "huanshi_token_expire"           //幻势:幻势:口令有效期(秒)
	DictionaryHuanshiSwitch               string = "huanshi_switch"                 //幻势:开关
	DictionaryHuanshiUrl                  string = "huanshi_url"                    //幻势:服务器地址
	DictionaryHuanshiAlarmSwitch          string = "huanshi_alarm_switch"           //幻势:告警开关
	DictionaryHuanshiAlarmIp              string = "huanshi_alarm_ip"               //幻势:告警ip
	DictionaryHuanshiAlarmPort            string = "huanshi_alarm_port"             //幻势:告警端口
	DictionaryHuanshiLogSwitch            string = "huanshi_log_switch"             //幻势:日志开关
	DictionaryHuanshiLogIp                string = "huanshi_log_ip"                 //幻势:日志ip
	DictionaryHuanshiLogPort              string = "huanshi_log_port"               //幻势:日志端口
	DictionaryHuanshiDefaultVisitor       string = "huanshi_default_visitor"        //幻势:默认访客
	DictionaryHuanshiAlarmSwitchProcess   string = "huanshi_alarm_switch_process"   //幻势:告警推送开关-进程
	DictionaryHuanshiAlarmSwitchDll       string = "huanshi_alarm_switch_dll"       //幻势:告警推送开关-模块
	DictionaryHuanshiAlarmSwitchTamper    string = "huanshi_alarm_switch_tamper"    //幻势:告警开关-防篡改
	DictionaryHuanshiAlarmSwitchBlackmail string = "huanshi_alarm_switch_blackmail" //幻势:告警开关-勒索
	DictionaryHuanshiAlarmSwitchPort      string = "huanshi_alarm_switch_port"      //幻势:告警开关-端口虚开
	DictionaryHuanshiAlarmSwitchUsb       string = "huanshi_alarm_switch_usb"       //幻势:告警开关-外设
	DictionaryHuanshiAlarmSwitchPlug      string = "huanshi_alarm_switch_plug"      //幻势:告警开关动态封堵
	DictionaryHuanshiAlarmSwitchLogin     string = "huanshi_alarm_switch_login"     //幻势:告警开关-登录口令破解
	DictionaryHuanshiAlarmSwitchOutreach  string = "huanshi_alarm_switch_outreach"  //幻势:告警开关-外联
	DictionaryHuanshiAlarmSwitchSelf      string = "huanshi_alarm_switch_self"      //幻势:告警开关-自保

	DictionarySystemDeviceId        string = "system_device_id"         //系统:设备标识码
	DictionarySystemPlatformUrl     string = "system_platform_url"      //系统:平台访问地址
	DictionarySystemPlatformIp      string = "system_platform_ip"       //系统:平台IP地址
	DictionarySystemPlatformPort    string = "system_platform_port"     //系统:平台开放端口
	DictionarySystemAgentUrl        string = "system_agent_url"         //系统:终端使用的访问地址
	DictionarySystemAgentIp         string = "system_agent_ip"          //系统:终端使用的IP地址
	DictionarySystemAgentPort       string = "system_agent_port"        //系统:终端使用的开放端口
	DictionarySystemGrantAutoSwitch string = "system_grant_auto_switch" //系统:自动授权开关
	DictionarySystemAlarmExpire     string = "system_alarm_expire"      //系统:告警失效时间(天)，最大180
	DictionarySystemSyslogExpire    string = "system_syslog_expire"     //系统:通信日志保留时长(天),最大180
	DictionarySystemSshSwitch       string = "system_ssh_switch"        //系统:ssh开关

)

var EnumDictionarySlice = map[string][]string{
	DictionaryModuleAdmin: {
		DictionaryAdminLoginExpire,
		DictionaryAdminPwLengthMin,
		DictionaryAdminPwLengthMax,
		DictionaryAdminLoginErrorMax,
		DictionaryAdminLoginErrorRemind,
		DictionaryAdminLoginLockDuration,
	},
	DictionaryModuleAgent: {
		DictionaryAgentLogSendProto,
		DictionaryAgentSaveAccessLog,
		DictionaryAgentTaskCacheExpire,
		DictionaryAgentGlobalPlugConfig,
		DictionaryAgentGlobalPlugLoginConfig,
		DictionaryAgentGlobalLoginSolveConfig,
	},
	DictionaryModuleThreaten: {
		DictionaryThreatenSwitch,
		DictionaryThreatenUrl,
		DictionaryThreatenAppId,
		DictionaryThreatenAppKey,
		DictionaryThreatenCacheExpire,
	},
	DictionaryModuleHuanshi: {
		DictionaryHuanshiTokenExpire,
		DictionaryHuanshiSwitch,
		DictionaryHuanshiUrl,
		DictionaryHuanshiAlarmSwitch,
		DictionaryHuanshiAlarmIp,
		DictionaryHuanshiAlarmPort,
		DictionaryHuanshiLogSwitch,
		DictionaryHuanshiLogIp,
		DictionaryHuanshiLogPort,
		DictionaryHuanshiDefaultVisitor,
		DictionaryHuanshiAlarmSwitchProcess,
		DictionaryHuanshiAlarmSwitchDll,
		DictionaryHuanshiAlarmSwitchTamper,
		DictionaryHuanshiAlarmSwitchBlackmail,
		DictionaryHuanshiAlarmSwitchPort,
		DictionaryHuanshiAlarmSwitchUsb,
		DictionaryHuanshiAlarmSwitchPlug,
		DictionaryHuanshiAlarmSwitchLogin,
		DictionaryHuanshiAlarmSwitchOutreach,
		DictionaryHuanshiAlarmSwitchSelf,
	},
	DictionaryModuleSystem: {
		DictionarySystemDeviceId,
		DictionarySystemPlatformUrl,
		DictionarySystemPlatformIp,
		DictionarySystemPlatformPort,
		DictionarySystemAgentUrl,
		DictionarySystemAgentIp,
		DictionarySystemAgentPort,
		DictionarySystemGrantAutoSwitch,
		DictionarySystemAlarmExpire,
		DictionarySystemSyslogExpire,
		DictionarySystemSshSwitch,
	},
}
