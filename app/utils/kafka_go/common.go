package kafka_go

type MessageStruct struct {
	Sign    string `json:"sign"`
	Type    string `json:"type"`
	OrgId   int64  `json:"oId"`
	AgentId string `json:"aId"`
	AgentOS string `json:"aOS"`
	Data    string `json:"data"`
}

const (
	TopicCommon   string = "TopicCommon"
	TopicProcess  string = "TopicProcess"
	TopicThreaten string = "TopicThreaten"
	TopicAlarm    string = "TopicAlarm"
	TopicSyslog   string = "TopicSyslog"

	TypeAgentDelete        string = "agent_delete"          //类型：终端删除
	TypeAgentPutMonitor    string = "agent_put_monitor"     //类型：终端上传监控数据
	TypeAgentOpenAutoWhite string = "agent_open_auto_white" //类型：终端开启自动白名单

	TypeOrgClean string = "org_clean" //类型： 企业清空业务数据

	TypeProcessPutMaster  string = "process_put_master"  //类型：上传主进程
	TypeProcessPutRunning string = "process_put_running" //类型：上传运行中进程
	TypeDllPutMaster      string = "dll_put_master"      //类型：上传主模块
	TypeDllPutRunning     string = "dll_put_running"     //类型：上传运行中模块

	TypeAlarmPutProcess   string = "alarm_put_process"   //类型：进程告警上传
	TypeAlarmPutDll       string = "alarm_put_dll"       //类型：模块告警上传
	TypeAlarmPutTamper    string = "alarm_put_tamper"    //类型：防篡改告警上传
	TypeAlarmPutBlackmail string = "alarm_put_blackmail" //类型：勒索告警上传
	TypeAlarmPutUsb       string = "alarm_put_usb"       //类型：外设告警上传
	TypeAlarmPutSelf      string = "alarm_put_self"      //类型：自保告警上传
	TypeAlarmPutPort      string = "alarm_put_port"      //类型：端口虚开告警上传

	TypeSyslogPutDns     string = "syslog_put_dns"     //类型：通讯日志-DNS上传
	TypeSyslogPutInner   string = "syslog_put_inner"   //类型：通讯日志-内联上传
	TypeSyslogPutOuter   string = "syslog_put_outer"   //类型：通讯日志-外联上传
	TypeSyslogPutProcess string = "syslog_put_process" //类型：通讯日志-进程启动上传
	TypeSyslogPutLogin   string = "syslog_put_login"   //类型：通讯日志-登录日志上传
)
