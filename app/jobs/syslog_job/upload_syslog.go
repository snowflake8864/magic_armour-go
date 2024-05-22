package syslog_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/agent/syslog_service/job_callback_syslog"
	"goskeleton/app/utils/helpers/json_helper"
)

func UploadDns(jsonStr string) {
	model := syslog_models.SyslogDnsModel{}
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_syslog.CreateDnsUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}

func UploadProcess(jsonStr string) {
	model := syslog_models.SyslogProcessModel{}
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_syslog.CreateProcessUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}
