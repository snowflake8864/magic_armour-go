package proc_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service/job_callback_proc"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
)

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
