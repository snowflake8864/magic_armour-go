package agent_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/agent/agent_service/agent_callback_service"
)

func Delete(orgId int64, agentId string) {
	variable.ZapLog.Info(fmt.Sprintf("orgId=%d, agentId=%s", orgId, agentId))
	agent_callback_service.CreateCurdFactory().DeleteByOrgAndAgent(orgId, agentId)
}
