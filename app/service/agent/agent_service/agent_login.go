package agent_service

import "goskeleton/app/model/model_agent"

func CreateLoginFactory() *AgentLogin {
	return &AgentLogin{model_agent.CreateAgentFactory()}
}

type AgentLogin struct {
	masterModel *model_agent.AgentModel
}
