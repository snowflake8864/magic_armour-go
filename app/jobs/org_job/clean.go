package org_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
)

func Clean(orgId int64) {
	variable.ZapLog.Info("org_job Clean" + helpers.JsonEncode(orgId))

	model_agent.CreateAgentFactory().DeleteByOrg(orgId)
	model_agent.CreateTaskFactory().DeleteByOrg(orgId)

}
