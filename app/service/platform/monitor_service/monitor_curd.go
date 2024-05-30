package monitor_service

import (
	"goskeleton/app/model/platform/server_monitor_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateCurdFactory() *MonitorCurd {
	return &MonitorCurd{masterModel: server_monitor_models.CreateServerMonitorFactory()}
}

type MonitorCurd struct {
	masterModel *server_monitor_models.ServerMonitorModel
}

// Add 新增
func (crd *MonitorCurd) Add(param *server_monitor_models.ServerMonitorModel) bool {
	if string_helper.IsEmpty(param.CreatedAt) {
		param.CreatedAt = date_helper.CurDatetime()
	}
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
