package enum_dictionary

// 数据字典

const (
	ModuleAdmin    string = "admin"    //模块：管路员
	ModuleAgent    string = "agent"    //模块：终端
	ModuleThreaten string = "threaten" //模块：终端
	ModuleHuanshi  string = "huanshi"  //模块：幻势
	ModuleSystem   string = "system"   //模块：系统

	AdminLoginExpire       string = "admin_login_expire"        //管理员:登录有效期(时)
	AdminPwLengthMin       string = "admin_pw_length_min"       //管理员:密码长度最小值(6)
	AdminPwLengthMax       string = "admin_pw_length_max"       //管理员:密码长度最小值(18)
	AdminLoginErrorMax     string = "admin_login_error_max"     //管理员:登录错误次数最大值(100)
	AdminLoginErrorRemind  string = "admin_login_error_remind"  //管理员:登录错误次数提醒
	AdminLoginLockDuration string = "admin_login_lock_duration" // 管理员:登录锁定时长(分)

	AgentLogSendProto           string = "agent_log_send_proto"            //终端:日志发送协议(1udp 2tcp)
	AgentSaveAccessLog          string = "agent_save_access_log"           //终端:任务缓存过期时间(时)
	AgentGlobalPlugPortConfig   string = "agent_global_plug_port_config"   //终端:全局封堵-端口虚开配置
	AgentGlobalPlugLoginConfig  string = "agent_global_plug_login_config"  //终端：全局封堵-登录破解配置
	AgentGlobalLoginSolveConfig string = "agent_global_Login_solve_config" //终端：全局登录口令暴力破解告警配置

	ThreatenSwitch      string = "threaten_switch"       //威胁情报:开关
	ThreatenUrl         string = "threaten_url"          //威胁情报:服务器地址
	ThreatenAppId       string = "threaten_app_id"       //威胁情报:应用id
	ThreatenAppKey      string = "threaten_app_key"      //威胁情报:应用key
	ThreatenCacheExpire string = "threaten_cache_expire" //威胁情报:缓存结果保留周期(天)

	HuanshiSwitch               string = "huanshi_switch"                 //幻势:开关
	HuanshiUrl                  string = "huanshi_url"                    //幻势:服务器地址
	HuanshiAlarmSwitch          string = "huanshi_alarm_switch"           //幻势:告警开关
	HuanshiAlarmIp              string = "huanshi_alarm_ip"               //幻势:告警ip
	HuanshiAlarmPort            string = "huanshi_alarm_port"             //幻势:告警端口
	HuanshiLogSwitch            string = "huanshi_log_switch"             //幻势:日志开关
	HuanshiLogIp                string = "huanshi_log_ip"                 //幻势:日志ip
	HuanshiLogPort              string = "huanshi_log_port"               //幻势:日志端口
	HuanshiDefaultVisitor       string = "huanshi_default_visitor"        //幻势:默认访客
	HuanshiAlarmSwitchProcess   string = "huanshi_alarm_switch_process"   //幻势:告警开关-进程
	HuanshiAlarmSwitchDll       string = "huanshi_alarm_switch_dll"       //幻势:告警开关-模块
	HuanshiAlarmSwitchTamper    string = "huanshi_alarm_switch_tamper"    //幻势:告警开关-防篡改
	HuanshiAlarmSwitchBlackmail string = "huanshi_alarm_switch_blackmail" //幻势:告警开关-勒索
	HuanshiAlarmSwitchPort      string = "huanshi_alarm_switch_port"      //幻势:告警开关-端口虚开
	HuanshiAlarmSwitchUsb       string = "huanshi_alarm_switch_usb"       //幻势:告警开关-外设
	HuanshiAlarmSwitchPlug      string = "huanshi_alarm_switch_plug"      //幻势:告警开关-动态封堵
	HuanshiAlarmSwitchLogin     string = "huanshi_alarm_switch_login"     //幻势:告警开关-登录口令破解
	HuanshiAlarmSwitchOutreach  string = "huanshi_alarm_switch_outreach"  //幻势:告警开关-外联
	HuanshiAlarmSwitchSelf      string = "huanshi_alarm_switch_self"      //幻势:告警开关-自保

	SystemDeviceId        string = "system_device_id"         //系统:设备标识码
	SystemPlatformProto   string = "system_platform_proto"    //系统:平台访问协议
	SystemPlatformIp      string = "system_platform_ip"       //系统:平台访问IP
	SystemPlatformPort    string = "system_platform_port"     //系统:平台访问端口
	SystemAgentProto      string = "system_agent_proto"       //系统:终端访问协议
	SystemAgentIp         string = "system_agent_ip"          //系统:终端访问IP
	SystemAgentPort       string = "system_agent_port"        //系统:终端使访问端口
	SystemGrantAutoSwitch string = "system_grant_auto_switch" //系统:自动授权开关
	SystemAlarmExpire     string = "system_alarm_expire"      //系统:告警失效时间(天)，最大180
	SystemSyslogExpire    string = "system_syslog_expire"     //系统:通信日志保留时长(天),最大180
	SystemSshSwitch       string = "system_ssh_switch"        //系统:ssh开关

)

var EnumSlice = map[string][]string{
	ModuleAdmin: {
		AdminLoginExpire,
		AdminPwLengthMin,
		AdminPwLengthMax,
		AdminLoginErrorMax,
		AdminLoginErrorRemind,
		AdminLoginLockDuration,
	},
	ModuleAgent: {
		AgentLogSendProto,
		AgentSaveAccessLog,
		AgentGlobalPlugPortConfig,
		AgentGlobalPlugLoginConfig,
		AgentGlobalLoginSolveConfig,
	},
	ModuleThreaten: {
		ThreatenSwitch,
		ThreatenUrl,
		ThreatenAppId,
		ThreatenAppKey,
		ThreatenCacheExpire,
	},
	ModuleHuanshi: {
		HuanshiSwitch,
		HuanshiUrl,
		HuanshiAlarmSwitch,
		HuanshiAlarmIp,
		HuanshiAlarmPort,
		HuanshiLogSwitch,
		HuanshiLogIp,
		HuanshiLogPort,
		HuanshiDefaultVisitor,
		HuanshiAlarmSwitchProcess,
		HuanshiAlarmSwitchDll,
		HuanshiAlarmSwitchTamper,
		HuanshiAlarmSwitchBlackmail,
		HuanshiAlarmSwitchPort,
		HuanshiAlarmSwitchUsb,
		HuanshiAlarmSwitchPlug,
		HuanshiAlarmSwitchLogin,
		HuanshiAlarmSwitchOutreach,
		HuanshiAlarmSwitchSelf,
	},
	ModuleSystem: {
		SystemDeviceId,
		SystemPlatformProto,
		SystemPlatformIp,
		SystemPlatformPort,
		SystemAgentProto,
		SystemAgentIp,
		SystemAgentPort,
		SystemGrantAutoSwitch,
		SystemAlarmExpire,
		SystemSyslogExpire,
		SystemSshSwitch,
	},
}
