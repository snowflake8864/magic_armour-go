package agent_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_agent"
)

func Delete(orgId int64, agentId string) {
	variable.ZapLog.Info(fmt.Sprintf("orgId=%d, agentId=%s", orgId, agentId))
	model_agent.CreateAgentFactory().DeleteByOrgAndAgent(orgId, agentId)
}
