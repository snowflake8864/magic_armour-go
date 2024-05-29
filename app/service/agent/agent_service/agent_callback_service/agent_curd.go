package agent_callback_service

import (
	"goskeleton/app/model/agent/agent_models"
)

func CreateCurdFactory() *AgentCurd {
	return &AgentCurd{masterModel: agent_models.CreateAgentFactory()}
}

type AgentCurd struct {
	masterModel *agent_models.AgentModel
}

// DeleteByOrg 通过企业删除
func (crd *AgentCurd) DeleteByOrg(orgId int64) bool {
	return crd.masterModel.DeleteByOrg(orgId)
}

// DeleteByOrgAndAgent 通过企业和终端删除
func (crd *AgentCurd) DeleteByOrgAndAgent(orgId int64, id string) bool {
	return crd.masterModel.DeleteByOrgAndAgent(orgId, id)
}
