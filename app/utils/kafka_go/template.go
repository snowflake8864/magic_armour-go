package kafka_go

type MessageStruct struct {
	Sign string
	Type string
	Data interface{}
}

const (
	TopicCommon   string = "TopicCommon"
	TopicProcess  string = "TopicProcess"
	TopicThreaten string = "TopicThreaten"
	TopicAlarm    string = "TopicAlarm"
	TopicSyslog   string = "TopicSyslog"

	SignAgent string = "agent"
	SignOrg   string = "org"
	SignProc  string = "proc"

	TypeAgentDelete string = "agent_delete"

	TypeOrgClean string = "org_clean"

	TypeProcMaster  string = "proc_master"
	TypeProcRunning string = "proc_running"
)
