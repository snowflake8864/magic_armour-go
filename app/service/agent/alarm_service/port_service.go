package alarm_service

import (
	"goskeleton/app/model/alarm/alarm_models"
)

func CreatePortCurdFactory() *AlarmPortCurd {
	return &AlarmPortCurd{
		alarmBase
		masterModel: alarm_models.CreatePortFactory(),
	}
}

type AlarmPortCurd struct {
	alarmBase
	masterModel *alarm_models.PortModel
}
