package agent_job

import (
	"fmt"
	"goskeleton/app/global/variable"
)

func Delete(agentId string) {
	variable.ZapLog.Info(fmt.Sprintf("agent_job Delete  agentId:%s", agentId))
}
