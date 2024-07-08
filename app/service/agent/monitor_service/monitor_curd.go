package monitor_service

import (
	"errors"
	"goskeleton/app/model/platform/server_monitor_models"
)

func CreateMonitorFactory() *MonitorCurd {
	return &MonitorCurd{server_monitor_models.CreateServerMonitorFactory()}
}

type MonitorCurd struct {
	masterModel *server_monitor_models.ServerMonitorModel
}

// Add 新增
func (crd *MonitorCurd) Add(param *server_monitor_models.ServerMonitorModel) error {

	if crd.masterModel.InsertItem(param) {
		return nil
	}

	return errors.New("保存失败")
}
