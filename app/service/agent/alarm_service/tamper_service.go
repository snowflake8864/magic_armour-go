package alarm_service

import (
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateTamperCurdFactory() *AlarmTamperCurd {
	return &AlarmTamperCurd{
		alarmBase
		masterModel: alarm_models.CreateTamperFactory(),
	}
}

type AlarmTamperCurd struct {
	alarmBase
	masterModel *alarm_models.TamperModel
}
