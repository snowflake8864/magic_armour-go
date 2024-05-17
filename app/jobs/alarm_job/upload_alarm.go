package alarm_job

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/alarm/alarm_models"
	"goskeleton/app/service/agent/alarm_service/job_callback_alarm"
	"goskeleton/app/utils/helpers/json_helper"
)

func UploadProcess(jsonStr string) {
	model := alarm_models.AlarmProcessModel{}
	if err := json_helper.Json2Struct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateProcessUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}

func UploadDll(jsonStr string) {
	model := alarm_models.AlarmDllModel{}
	if err := json_helper.Json2Struct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateDllUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}

func UploadTamper(jsonStr string) {
	model := alarm_models.AlarmTamperModel{}
	if err := json_helper.Json2Struct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateTamperUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}

func UploadBlackmail(jsonStr string) {
	model := alarm_models.AlarmBlackmailModel{}
	if err := json_helper.Json2Struct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateBlackmailUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}

func UploadUsb(jsonStr string) {
	model := alarm_models.AlarmUsbModel{}
	if err := json_helper.Json2Struct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateUsbUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}

func UploadSelf(jsonStr string) {
	model := alarm_models.AlarmUsbModel{}
	if err := json_helper.Json2Struct(jsonStr, &model); err != nil {
		variable.ZapLog.Info(err.Error())
		variable.ZapLog.Info(jsonStr)
		return
	}

	if err := job_callback_alarm.CreateUsbUploadFactory().Add(&model); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}
}
