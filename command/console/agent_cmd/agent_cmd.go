package agent_cmd

import "goskeleton/app/global/variable"

func CreateCmdFactory() *AgentCmd {
	return &AgentCmd{}
}

type AgentCmd struct {
}

// Offline 离线
func (crd *AgentCmd) Offline() {
	variable.ZapLog.Sugar().Debug("Offline")
}

// HandleSafe 处理安全状态
func (crd *AgentCmd) HandleSafe() {
	variable.ZapLog.Sugar().Debug("HandleSafe")
}

// GrantExpired  授权已过期的
func (crd *AgentCmd) GrantExpired() {
	variable.ZapLog.Sugar().Debug("GrantExpired")
}
