package proc_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service/job_callback_proc"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
)

func UploadProcess(jsonStr string) {
	processModel := proc_models.ProcessModel{}
	if err := json_helper.ToStruct(jsonStr, &processModel); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_proc.CreateProcessUploadFactory().Add(&processModel); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func BatchUploadProcess(orgId int64, agentId string, jsonStr string) {
	fmt.Println("BatchUploadProcess start: ", date_helper.CurTimestamp())
	var models []proc_models.ProcessModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_proc.CreateProcessUploadFactory().BatchAdd(orgId, agentId, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
	fmt.Println("BatchUploadProcess end: ", date_helper.CurTimestamp())
}

func UploadProcessRunning(jsonStr string) {
	processRunningModel := proc_models.ProcessRunningModel{}
	if err := json_helper.ToStruct(jsonStr, &processRunningModel); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_proc.CreateProcessRunningUploadFactory().Add(&processRunningModel); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}
