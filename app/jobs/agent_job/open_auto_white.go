package agent_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/proc_service"
)

func OpenAutoWhite(orgId int64, agentId string) {
	variable.ZapLog.Debug(fmt.Sprintf("orgId=%d, agentId=%s", orgId, agentId))

	AgentFact := agent_models.CreateAgentFactory()
	if agentModel := AgentFact.GetItemByOrgAndAgent(orgId, agentId, AgentFact.ColumnProc(true)); agentModel != nil && agentModel.OrgId > 0 {
		proc_service.CreateProcessCurdFactory().AutoWhiteSwitchOnAgent(orgId, agentId, agentModel.AutoWhiteSwitch)
		proc_service.CreateDllCurdFactory().AutoWhiteSwitchOnAgent(orgId, agentId, agentModel.AutoWhiteSwitch)
	}

}
