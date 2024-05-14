package agent_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_model"
)

func Delete(orgId int64, agentId string) {
	variable.ZapLog.Info(fmt.Sprintf("orgId=%d, agentId=%s", orgId, agentId))
	agent_model.CreateAgentFactory().DeleteByOrgAndAgent(orgId, agentId)
}
