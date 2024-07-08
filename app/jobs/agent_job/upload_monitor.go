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
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := agent_monitor_service.CreateMonitorFactory().Add(&model); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}
