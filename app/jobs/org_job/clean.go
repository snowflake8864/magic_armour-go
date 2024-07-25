package org_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/task_models"
	"goskeleton/app/service/agent/agent_service/agent_callback_service"
)

func Clean(orgId int64) {
	variable.ZapLog.Sugar().Info(fmt.Sprintf("orgId=%d", orgId))

	agent_callback_service.CreateCurdFactory().DeleteByOrg(orgId)
	task_models.CreateTaskFactory().DeleteByOrg(orgId)

}
