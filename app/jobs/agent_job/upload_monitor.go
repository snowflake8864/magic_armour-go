package agent_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/service/agent/agent_monitor_service"
	"goskeleton/app/utils/helpers/json_helper"
)

func UploadMonitor(jsonStr string) {
	model := agent_monitor_models.MonitorModel{}
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Error(err.Error())
		variable.ZapLog.Error(jsonStr)
		return
	}

	if err := agent_monitor_service.CreateMonitorFactory().Add(&model); err != nil {
		variable.ZapLog.Error(err.Error())
		return
	}
}
