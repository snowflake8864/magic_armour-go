package org_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/task_models"
)

func Clean(orgId int64) {
	variable.ZapLog.Info(fmt.Sprintf("orgId=%d", orgId))

	agent_models.CreateAgentFactory().DeleteByOrg(orgId)
	task_models.CreateTaskFactory().DeleteByOrg(orgId)

}
