package syslog_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/agent/syslog_service/job_callback_syslog"
	"goskeleton/app/utils/helpers/json_helper"
)

func UploadDns(jsonStr string) {
	model := syslog_models.DnsModel{}
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_syslog.CreateDnsUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func UploadOuter(jsonStr string) {
	var models []syslog_models.OuterModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_syslog.CreateOuterUploadFactory().BatchAdd(&models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func UploadInner(jsonStr string) {
	var models []syslog_models.InnerModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_syslog.CreateInnerUploadFactory().BatchAdd(&models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func UploadProcess(jsonStr string) {
	var models []syslog_models.ProcessModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_syslog.CreateProcessUploadFactory().BatchAdd(&models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func UploadLogin(jsonStr string) {
	model := syslog_models.LoginModel{}
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_syslog.CreateLoginUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}
