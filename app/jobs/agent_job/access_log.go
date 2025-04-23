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
		variable.ZapLog.Warn(err.Error())
		variable.ZapLog.Warn(jsonStr)
		return
	}

	if !access_log_service.CreateCurdFactory().Add(&model) {
		variable.ZapLog.Warn("Add Fail" + jsonStr)
		return
	}
}
