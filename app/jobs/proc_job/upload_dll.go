package proc_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service/job_callback_proc"
	"goskeleton/app/utils/helpers/json_helper"
)

func UploadDll(jsonStr string) {
	dllModel := proc_models.DllModel{}
	if err := json_helper.Json2Struct(jsonStr, &dllModel); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_proc.CreateDllUploadFactory().Add(&dllModel); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}

func UploadDllRunning(jsonStr string) {
	dllRunningModel := proc_models.DllRunningModel{}
	if err := json_helper.Json2Struct(jsonStr, &dllRunningModel); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_proc.CreateDllRunningUploadFactory().Add(&dllRunningModel); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}
