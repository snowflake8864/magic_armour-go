package kafka_go

type MessageStruct struct {
	Sign string
	Type string
	Data interface{}
}

const (
	SignAgent       string = "agent"
	TypeAgentDelete string = "agent_delete"
)
