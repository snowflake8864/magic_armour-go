package alarm_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/agent/alarm_service/job_callback_alarm"
	"goskeleton/app/utils/helpers/json_helper"
)

func UploadProcess(jsonStr string) {
	var models []alarm_models.ProcessModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateProcessUploadFactory().BatchAdd(&models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func UploadDll(jsonStr string) {
	var models []alarm_models.DllModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}
	if err := job_callback_alarm.CreateDllUploadFactory().BatchAdd(&models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func UploadTamper(jsonStr string) {
	var models []alarm_models.TamperModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateTamperUploadFactory().BatchAdd(&models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func UploadBlackmail(jsonStr string) {
	var models []alarm_models.BlackmailModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateBlackmailUploadFactory().BatchAdd(&models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func UploadUsb(jsonStr string) {
	var models []alarm_models.UsbModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateUsbUploadFactory().BatchAdd(&models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func UploadSelf(jsonStr string) {
	var models []alarm_models.SelfModel
	if err := json_helper.ToStruct(jsonStr, &models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateSelfUploadFactory().BatchAdd(&models); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}

func UploadPort(jsonStr string) {
	model := alarm_models.PortModel{}
	if err := json_helper.ToStruct(jsonStr, &model); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		variable.ZapLog.Sugar().Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreatePortUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
		return
	}
}
