package alarm_service

import (
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateSelfCurdFactory() *AlarmSelfCurd {
	return &AlarmSelfCurd{
		alarmBase
		masterModel: alarm_models.CreateSelfFactory(),
	}
}

type AlarmSelfCurd struct {
	alarmBase
	masterModel *alarm_models.SelfModel
}
