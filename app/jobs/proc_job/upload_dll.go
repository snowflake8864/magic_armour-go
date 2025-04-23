package proc_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service/job_callback_proc"
	"goskeleton/app/utils/helpers/json_helper"
)

func BatchUploadDll(orgId int64, agentId string, jsonStr string) {
	var models []proc_models.DllModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Warn(err.Error())
		variable.ZapLog.Warn(jsonStr)
		return
	}

	if err := job_callback_proc.CreateDllUploadFactory().BatchAdd(orgId, agentId, &models); err != nil {
		variable.ZapLog.Warn(err.Error())
		return
	}
}
