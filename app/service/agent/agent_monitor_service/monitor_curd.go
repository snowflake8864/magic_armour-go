package agent_monitor_service

import (
	"errors"
	"goskeleton/app/model/agent/agent_monitor_models"
)

func CreateMonitorFactory() *MonitorCurd {
	return &MonitorCurd{agent_monitor_models.CreateMonitorFactory()}
}

type MonitorCurd struct {
	masterModel *agent_monitor_models.MonitorModel
}

// Add 新增
func (crd *MonitorCurd) Add(param *agent_monitor_models.MonitorModel) error {

	if crd.masterModel.InsertItem(param) {
		return nil
	}

	return errors.New("保存失败")
}

func (crd *MonitorCurd) sendEmail(model *agent_monitor_models.MonitorModel) {

}
