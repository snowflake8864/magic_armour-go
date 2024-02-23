package agent_service

import (
	"goskeleton/app/model/model_agent"
)

func CreateCurdFactory() *AgentCurd {
	return &AgentCurd{model_agent.CreateAgentFactory()}
}

type AgentCurd struct {
	masterModel *model_agent.AgentModel
}
