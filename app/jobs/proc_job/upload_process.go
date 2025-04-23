package proc_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service/job_callback_proc"
	"goskeleton/app/utils/helpers/json_helper"
)

func BatchUploadProcess(orgId int64, agentId string, jsonStr string) {
	var models []proc_models.ProcessModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Debug(err.Error())
		variable.ZapLog.Debug(jsonStr)
		return
	}

	if err := job_callback_proc.CreateProcessUploadFactory().BatchAdd(orgId, agentId, &models); err != nil {
		variable.ZapLog.Debug(err.Error())
		return
	}
}
