package syslog_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/agent/alarm_service/job_callback_alarm"
	"goskeleton/app/utils/helpers/json_helper"
)

func UploadProcess(jsonStr string) {
	model := syslog_models.CreateSyslogProcessFactory()
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateProcessUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}
