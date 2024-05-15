package proc_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service/job_callcak_proc"
	"goskeleton/app/utils/helpers"
)

func UploadProcess(jsonStr string) {
	processModel := proc_models.ProcessModel{}
	if err := helpers.Json2Struct(jsonStr, &processModel); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callcak_proc.CreateProcessUploadFactory().Add(&processModel); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}

func UploadProcessRunning(jsonStr string) {
	processRunningModel := proc_models.ProcessRunningModel{}
	if err := helpers.Json2Struct(jsonStr, &processRunningModel); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callcak_proc.CreateProcessRunningUploadFactory().Add(&processRunningModel); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}
