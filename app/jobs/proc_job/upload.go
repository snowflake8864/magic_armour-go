package proc_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
)

func UploadProcess(model *model_agent.ProcessModel) {
	variable.ZapLog.Info("proc_job UploadProcess" + helpers.JsonEncode(model))
}

func UploadProcessRunning(model *model_agent.ProcessRunningModel) {
	variable.ZapLog.Info("proc_job UploadProcessRunning" + helpers.JsonEncode(model))
}
