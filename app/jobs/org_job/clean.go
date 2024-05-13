package org_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_agent"
)

func Clean(orgId int64) {
	variable.ZapLog.Info(fmt.Sprintf("orgId=%d", orgId))

	model_agent.CreateAgentFactory().DeleteByOrg(orgId)
	model_agent.CreateTaskFactory().DeleteByOrg(orgId)

}
