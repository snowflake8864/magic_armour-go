package alarm_service

import (
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateProcessCurdFactory() *AlarmProcessCurd {
	return &AlarmProcessCurd{
		alarmBase: alarm_models.CreateProcessFactory(),
	}
}

type AlarmProcessCurd struct {
	alarmBase *alarm_models.ProcessModel
}
