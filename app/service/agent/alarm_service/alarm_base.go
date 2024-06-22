package alarm_service

import (
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/utils/helpers/date_helper"
)

type alarmBase struct {
}

func (crd *alarmBase) getUpdateForHandle(remark string, user string) map[string]interface{} {
	var updates = make(map[string]interface{})
	updates["handle_status"] = enum_alarm.HandleStatusYes
	updates["handle_remark"] = remark
	updates["handle_user"] = user
	updates["handle_at"] = date_helper.CurDatetime()
	return updates
}
