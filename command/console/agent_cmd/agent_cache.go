package agent_cmd

func CreateCmdFactory() *AgentCmd {
	return &AgentCmd{}
}

type AgentCmd struct {
}

// Offline 离线
func (crd *AgentCmd) Offline() {

}

// IsSafe 是否是否安全
func (crd *AgentCmd) IsSafe() {

}
