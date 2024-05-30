package agent_job

import (
	"fmt"
	"goskeleton/app/global/variable"
)

func OpenAutoWhite(orgId int64, agentId string) {
	variable.ZapLog.Info(fmt.Sprintf("orgId=%d, agentId=%s", orgId, agentId))
}
