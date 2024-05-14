package org_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_model"
	"goskeleton/app/model/model_agent"
)

func Clean(orgId int64) {
	variable.ZapLog.Info(fmt.Sprintf("orgId=%d", orgId))

	agent_model.CreateAgentFactory().DeleteByOrg(orgId)
	model_agent.CreateTaskFactory().DeleteByOrg(orgId)

}
