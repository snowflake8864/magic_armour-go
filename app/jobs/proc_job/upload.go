package proc_job

import (
	"goskeleton/app/global/variable"
)

func UploadProcess(jsonStr string) {
	variable.ZapLog.Info("proc_job UploadProcess" + jsonStr)
}

func UploadProcessRunning(jsonStr string) {
	variable.ZapLog.Info("proc_job UploadProcessRunning" + jsonStr)
}
