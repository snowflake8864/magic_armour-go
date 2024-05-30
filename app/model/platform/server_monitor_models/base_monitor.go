package server_monitor_models

import "goskeleton/app/utils/helpers/string_helper"

type SearchConditions struct {
	MonitorDay string `json:"monitor_day"`
	ServerCode string `json:"server_code"`
}
type baseMonitor struct {
}

// search2Where 查询条件转化
func (crd *baseMonitor) conditions2Where(search SearchConditions) map[string]interface{} {
	wheres := make(map[string]interface{})
	if !string_helper.IsEmpty(search.MonitorDay) {
		wheres["monitor_day"] = search.MonitorDay
	}
	if !string_helper.IsEmpty(search.ServerCode) {
		wheres["ServerCode"] = search.ServerCode
	}
	return wheres
}
