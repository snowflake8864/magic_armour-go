package enum_dictionary

// 数据字典

const (
	ModulePlatform     string = "platform"     //模块：平台
	ModuleOrganization string = "organization" //模块：系统
	ModuleAgent        string = "agent"        //模块：终端
	ModuleThreaten     string = "threaten"     //模块：威胁情报
	ModuleHuanshi      string = "huanshi"      //模块：幻势
	ModuleOrder        string = "order"        //模块：订单

	PlatformLoginErrorMax     string = "platform_login_error_max"     //平台:登录错误次数最大值(100)
	PlatformLoginErrorRemind  string = "platform_login_error_remind"  //平台:登录错误次数提醒
	PlatformLoginLockDuration string = "platform_login_lock_duration" //平台:登录锁定时长(分)
	PlatformSshSwitch         string = "platform_ssh_switch"          //平台:ssh开关
	PlatformSshCmdStart       string = "platform_ssh_start"           //平台:ssh服务开启命令
	PlatformSshCmdStop        string = "platform_ssh_stop"            //平台:ssh服务停止命令
	PlatformServiceCodeList   string = "platform_service_code_list"   //平台:服务器编号列表
	PlatformTimeCmd           string = "platform_time_cmd"            //平台:时间修改命令
	PlatformInstallAt         string = "platform_install_at"          //平台:安装时间
	PlatformSntpAddr          string = "platform_sntp_addr"           //平台:SNTP服务器地址
	PlatformSntpCmd           string = "platform_sntp_cmd"            //平台:SNTP同步命令

	OrganizationLoginSmsSwitch       string = "org_login_sms_switch"       //企业:登录短信验证开关
	OrganizationLoginEXPIRE          string = "org_login_expire"           //企业:登录有效期(分)
	OrganizationLoginErrorMax        string = "org_login_error_max"        //企业:登录错误次数最大值(100)
	OrganizationLoginErrorRemind     string = "org_login_error_remind"     //企业:登录错误次数提醒
	OrganizationLoginLockDuration    string = "org_login_lock_duration"    //企业:登录锁定时长(分)
	OrganizationPwLengthMin          string = "org_pw_length_min"          //企业:密码长度最小值(6)
	OrganizationPwLengthMax          string = "org_pw_length_max"          //企业:密码长度最大值(18)
	OrganizationPwExpire             string = "org_pw_expire"              //企业:密码有效期(天)
	OrganizationPwCycleModifySwitch  string = "org_pw_cycle_modify_switch" //企业:密码周期修改开关
	OrganizationPwlFirstModifySwitch string = "org_pw_first_modify_switch" //企业::密码首次登录修改开关
	OrganizationSubUserLimit         string = "org_sub_user_limit"         //企业:子用户限制数目

	AgentHttpUrl               string = "agent_http_url"                //终端:http请求地址(无端口)
	AgentHttpPort              string = "agent_http_port"               //终端:http请求端口
	AgentLogSendProto          string = "agent_log_send_proto"          //终端:日志发送协议(1udp 2tcp)
	AgentSaveAccessLog         string = "agent_save_access_log"         //终端:是否保存访问记录(0|1)
	AgentConfigExpire          string = "agent_config_expire"           //终端:配置文件有效期(小时)
	AgentAutoDeleteExpire      string = "agent_auto_delete_expire"      //终端:自动删除有效期(天)
	AgentAlarmExpire           string = "agent_alarm_expire"            //终端:告警保留时长(天)
	AgentSyslogExpire          string = "agent_syslog_expire"           //终端:通讯日志保留时长(天)
	AgentProcGlobalSwitch      string = "agent_proc_global_switch"      //终端:进程全局名单是否下发给终端开关
	AgentProcGlobalLimit       string = "agent_proc_global_limit"       //终端:进程全局名单限制数目
	AgentCompatibleDefaultOrg  string = "agent_compatible_default_org"  //终端:兼容以前离线版本默认企业编号
	AgentCompatibleDefaultUser string = "agent_compatible_default_user" //终端:兼容以前离线版本默认用户编号
	AgentTypes                 string = "agent_types"                   //终端:终端类型-离线版 todo http://192.168.19.9/zentao/story-view-781.html
	AgentOutreachDetectMax     string = "agent_outreach_detect_max"     //终端:外联探测最大数
	AgentOutreachDetectTime    string = "agent_outreach_detect_time"    //终端:外联探测间隔时间(秒)
	AgentDebugIds              string = "agent_debug_ids"               //终端:调试终端集合
	AgentDebugPort             string = "agent_debug_port"              //终端:调试端口

	ThreatenSwitch        string = "threaten_switch"         //威胁情报:开关
	ThreatenUrl           string = "threaten_url"            //威胁情报:服务器地址
	ThreatenAppId         string = "threaten_app_id"         //威胁情报:应用id
	ThreatenAppKey        string = "threaten_app_key"        //威胁情报:应用key
	ThreatenAppExpire     string = "threaten_app_expire"     //威胁情报:应用过期时间
	ThreatenCacheExpire   string = "threaten_cache_expire"   //威胁情报:缓存结果保留周期(天)
	ThreatenSupportNumber string = "threaten_support_number" //威胁情报:校验支持数量,最大300

	HuanshiCallbackUrl string = "huanshi_callback_url" //幻势:回调地址(有端口)
	HuanshiCallIp      string = "huanshi_callback_ip"  //幻势:回调ip
	HuanshiSwitch      string = "huanshi_switch"       //幻势:开关
	HuanshiUrl         string = "huanshi_url"          //幻势:服务器地址
	HuanshiAlarmSwitch string = "huanshi_alarm_switch" //幻势:告警开关
	HuanshiAlarmIp     string = "huanshi_alarm_ip"     //幻势:告警ip
	HuanshiAlarmPort   string = "huanshi_alarm_port"   //幻势:告警端口
	HuanshiLogSwitch   string = "huanshi_log_switch"   //幻势:日志开关
	HuanshiLogIp       string = "huanshi_log_ip"       //幻势:日志ip
	HuanshiLogPort     string = "huanshi_log_port"     //幻势:日志端口
	HuanshiThirdSwitch string = "huanshi_third_switch" //幻势:友商开关
	HuanshiThirdIp     string = "huanshi_third_ip"     //幻势:友商ip
	HuanshiThirdPort   string = "huanshi_third_port"   //幻势:友商端口

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
	HuanshiLogSwitchOrgOperate  string = "huanshi_log_switch_org_operate" //幻势:日志开关-企业操作
	HuanshiLogSwitchDns         string = "huanshi_log_switch_dns"         //幻势:日志开关-DNS
	HuanshiLogSwitchInner       string = "huanshi_log_switch_inner"       //幻势:日志开关-对内通信
	HuanshiLogSwitchOuter       string = "huanshi_log_switch_outer"       //幻势:日志开关-对外通信
	HuanshiLogSwitchProcess     string = "huanshi_log_switch_process"     //幻势:日志开关-进程启动
	HuanshiLogSwitchLogin       string = "huanshi_log_switch_login"       //幻势:日志开关-登录

	OrderBankcardNo   string = "order_bankcard_no"   //订单：银行卡-卡号
	OrderBankcardName string = "order_bankcard_name" //订单：银行卡-卡名
	OrderBankcardBank string = "order_bankcard_bank" //订单：银行卡-开户行
	OrderPayingTel    string = "order_paying_tel"    //订单：支付时联系电话
	OrderExpireDays   string = "order_expire_days"   //订单：过期天数
	OrderLicenseName  string = "order_license_name"  //订单：导入证书默认名称
)

var EnumSlice = map[string][]string{
	ModulePlatform: {
		PlatformLoginErrorMax,
		PlatformLoginErrorRemind,
		PlatformLoginLockDuration,
		PlatformSshSwitch,
		PlatformSshCmdStart,
		PlatformSshCmdStop,
		PlatformServiceCodeList,
		PlatformTimeCmd,
		PlatformInstallAt,
		PlatformSntpAddr,
		PlatformSntpCmd,
	},
	ModuleOrganization: {
		OrganizationLoginSmsSwitch,
		OrganizationLoginEXPIRE,
		OrganizationLoginErrorMax,
		OrganizationLoginErrorRemind,
		OrganizationLoginLockDuration,
		OrganizationPwLengthMin,
		OrganizationPwLengthMax,
		OrganizationPwExpire,
		OrganizationPwCycleModifySwitch,
		OrganizationPwlFirstModifySwitch,
		OrganizationSubUserLimit,
	},
	ModuleAgent: {
		AgentHttpUrl,
		AgentHttpPort,
		AgentLogSendProto,
		AgentSaveAccessLog,
		AgentConfigExpire,
		AgentAutoDeleteExpire,
		AgentAlarmExpire,
		AgentSyslogExpire,
		AgentProcGlobalSwitch,
		AgentProcGlobalLimit,
		AgentCompatibleDefaultOrg,
		AgentCompatibleDefaultUser,
		AgentTypes,
		AgentOutreachDetectMax,
		AgentOutreachDetectTime,
		AgentDebugIds,
		AgentDebugPort,
	},

	ModuleThreaten: {
		ThreatenSwitch,
		ThreatenUrl,
		ThreatenAppId,
		ThreatenAppKey,
		ThreatenAppExpire,
		ThreatenCacheExpire,
		ThreatenSupportNumber,
	},
	ModuleHuanshi: {
		HuanshiCallbackUrl,
		HuanshiCallIp,
		HuanshiSwitch,
		HuanshiUrl,
		HuanshiAlarmSwitch,
		HuanshiAlarmIp,
		HuanshiAlarmPort,
		HuanshiLogSwitch,
		HuanshiLogIp,
		HuanshiLogPort,
		HuanshiThirdSwitch,
		HuanshiThirdIp,
		HuanshiThirdPort,
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
		HuanshiLogSwitchOrgOperate,
		HuanshiLogSwitchDns,
		HuanshiLogSwitchInner,
		HuanshiLogSwitchOuter,
		HuanshiLogSwitchProcess,
		HuanshiLogSwitchLogin,
	},
	ModuleOrder: {
		OrderBankcardNo,
		OrderBankcardName,
		OrderBankcardBank,
		OrderPayingTel,
		OrderExpireDays,
		OrderLicenseName,
	},
}
