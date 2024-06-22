package alarm_service

import (
	"goskeleton/app/model/alarm/alarm_models"
)

func CreatPlugCurdCurdFactory() *AlarmPlugCurd {
	return &AlarmPlugCurd{
		alarmBase
		masterModel: alarm_models.CreatePlugFactory(),
	}
}

type AlarmPlugCurd struct {
	alarmBase
	masterModel *alarm_models.PlugModel
}
