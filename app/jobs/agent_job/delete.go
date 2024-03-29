package agent_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/agent/agent_service"
)

func Delete(agentId string) {
	variable.ZapLog.Info(fmt.Sprintf("agent_job Delete  agentId:%s", agentId))
	agent_service.CreateCurdFactory().DeleteCallback(agentId)

}
