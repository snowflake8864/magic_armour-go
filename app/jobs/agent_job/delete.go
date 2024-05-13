package agent_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_agent"
)

func Delete(orgId int64, agentId string) {
	variable.ZapLog.Info(fmt.Sprintf("agent_job Delete orgId=%d, agentId=%s", orgId, agentId))
	model_agent.CreateAgentFactory().DeleteByOrgAndAgent(orgId, agentId)
}
