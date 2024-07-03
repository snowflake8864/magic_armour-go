package agent_cmd

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateMonitorFactory() *MonitorCmd {
	return &MonitorCmd{}
}

type MonitorCmd struct {
}

// Delete  删除
func (crd *MonitorCmd) Delete() {
	variable.ZapLog.Sugar().Debug("Delete start")
	_, _, dayEnd, _ := date_helper.CalcByDay(-90)
	agent_monitor_models.CreateMonitorFactory().DeleteByConditions(agent_monitor_models.SearchConditions{
		CreateAtLte: dayEnd,
	})
	variable.ZapLog.Sugar().Debug("Delete end")
}
