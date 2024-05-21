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

	TypeAgentDelete string = "agent_delete"

	TypeOrgClean string = "org_clean"

	TypeProcessPutMaster  string = "process_put_master"
	TypeProcessPutRunning string = "process_put_running"
	TypeDllPutMaster      string = "dll_put_master"
	TypeDllPutRunning     string = "dll_put_running"

	TypeAlarmPutProcess   string = "alarm_put_process"
	TypeAlarmPutDll       string = "alarm_put_dll"
	TypeAlarmPutTamper    string = "alarm_put_tamper"
	TypeAlarmPutBlackmail string = "alarm_put_blackmail"
	TypeAlarmPutUsb       string = "alarm_put_usb"
	TypeAlarmPutSelf      string = "alarm_put_self"
	TypeAlarmPutPort      string = "alarm_put_port"

	TypeSyslogPutDns     string = "syslog_put_dns"
	TypeSyslogPutInner   string = "syslog_put_inner"
	TypeSyslogPutOuter   string = "syslog_put_outer"
	TypeSyslogPutProcess string = "syslog_put_process"
	TypeSyslogPutLogin   string = "syslog_put_login"
)
