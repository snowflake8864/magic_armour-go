package proc_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/service/agent/proc_service/job_callcak_proc"
	"goskeleton/app/utils/helpers"
)

func UploadProcess(jsonStr string) {

	processModel := model_agent.ProcessModel{}
	if err := helpers.Json2Struct(jsonStr, &processModel); err != nil {
		variable.ZapLog.Info(jsonStr)
		variable.ZapLog.Info(err.Error())
		return
	}

	/*agentFact := agent_service.CreateCurdFactory()
	agentModel := agentFact.GetItemByOrgAndAgent(processModel.OrgId, processModel.AgentId, agentFact.ColumnProc(true))
	if agentModel == nil || helpers.StringIsEmpty(agentModel.AgentId) {
		variable.ZapLog.Info("终端不存在")
		return
	}*/

	if _, err := job_callcak_proc.CreateProcessUploadFactory().Add(&processModel); err != nil {
		variable.ZapLog.Info(jsonStr)
		variable.ZapLog.Info(err.Error())
		return
	}
}

func UploadProcessRunning(jsonStr string) {
	variable.ZapLog.Info(jsonStr)
}
