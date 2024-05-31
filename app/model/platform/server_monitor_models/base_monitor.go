package server_monitor_models

import (
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	MonitorDay string `json:"monitor_day"`
	ServerCode string `json:"server_code"`
}
type baseMonitor struct {
}

// search2Where 查询条件转化
func (crd *baseMonitor) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {

	if !string_helper.IsEmpty(condition.MonitorDay) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "monitor_day",
			Val: condition.MonitorDay,
		})
	}
	if !string_helper.IsEmpty(condition.ServerCode) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "server_code",
			Val: condition.ServerCode,
		})
	}
	return wheres
}
