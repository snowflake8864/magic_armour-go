package agent_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_access_models"
	"goskeleton/app/service/agent/access_log_service"
	"goskeleton/app/utils/helpers/json_helper"
)

func SaveAccessLog(jsonStr string) {
	model := agent_access_models.AccessLogModel{}
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Error(err.Error())
		variable.ZapLog.Debug(jsonStr)
		return
	}

	if !access_log_service.CreateCurdFactory().Add(&model) {
		variable.ZapLog.Error("Add Fail" + jsonStr)
		return
	}
}
