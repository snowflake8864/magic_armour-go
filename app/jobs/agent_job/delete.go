package agent_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
)

func Delete(param AgentCommon) {
	variable.ZapLog.Info("agent_job Delete" + helpers.JsonEncode(param))
	model_agent.CreateAgentFactory().DeleteByOrgAndAgent(param.OrgId, param.AgentId)
}
