package proc_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service/job_callback_proc"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
)

func UploadDll(jsonStr string) {
	dllModel := proc_models.DllModel{}
	if err := json_helper.ToStruct(jsonStr, &dllModel); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_proc.CreateDllUploadFactory().Add(&dllModel); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func BatchUploadDll(orgId int64, agentId string, jsonStr string) {
	fmt.Println("BatchUploadDll start: ", date_helper.CurTimestamp())
	var models []proc_models.DllModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_proc.CreateDllUploadFactory().BatchAdd(orgId, agentId, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
	fmt.Println("BatchUploadDll end: ", date_helper.CurTimestamp())
}

func UploadDllRunning(jsonStr string) {
	dllRunningModel := proc_models.DllRunningModel{}
	if err := json_helper.ToStruct(jsonStr, &dllRunningModel); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_proc.CreateDllRunningUploadFactory().Add(&dllRunningModel); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}
