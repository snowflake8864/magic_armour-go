package alarm_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/agent/proc_service/job_callback_proc"
	"goskeleton/app/utils/helpers/json_helper"
)

func UploadProcess(jsonStr string) {
	model := alarm_models.AlarmProcessModel{}
	if err := json_helper.Json2Struct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	variable.ZapLog.Info(model.ProcDir)
	if err := job_callback_proc.CreateDllUploadFactory().Add(&dllModel); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}
