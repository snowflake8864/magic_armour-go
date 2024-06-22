package alarm_service

import (
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateUsbCurdFactory() *AlarmUsbCurd {
	return &AlarmUsbCurd{
		alarmBase
		masterModel: alarm_models.CreateUsbFactory(),
	}
}

type AlarmUsbCurd struct {
	alarmBase
	masterModel *alarm_models.UsbModel
}
