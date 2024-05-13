package proc_job

import (
	"goskeleton/app/global/variable"
)

func UploadProcess(jsonStr string) {
	variable.ZapLog.Info(jsonStr)

}

func UploadProcessRunning(jsonStr string) {
	variable.ZapLog.Info(jsonStr)
}
