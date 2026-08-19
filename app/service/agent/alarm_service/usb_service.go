package alarm_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateUsbCurdFactory() *UsbCurd {
	return &UsbCurd{
		masterModel: alarm_models.CreateUsbFactory(),
	}
}

type UsbCurd struct {
	alarmBase
	masterModel *alarm_models.UsbModel
}

// Handle 处理
func (crd *UsbCurd) Handle(param *alarm_models.UsbModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:        param.OrgId,
		AgentId:      param.AgentId,
		UsbHash:      param.UsbHash,
		HandleStatus: enum_alarm.HandleStatusNot,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	return nil
}

// HandleAll 处理全部
func (crd *UsbCurd) HandleAll(param *alarm_models.UsbModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:        param.OrgId,
		AgentIds:     agentIds,
		HandleStatus: enum_alarm.HandleStatusNot,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New(my_errors.CurdSaveFailed)
	}
	return nil
}

// ClearByAgent 清空终端下告警
func (crd *UsbCurd) ClearByAgent(param *alarm_models.UsbModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New(my_errors.CurdClearFailed)
	}
	return nil
}

// ClearByAgents 清空多终端下告警
func (crd *UsbCurd) ClearByAgents(param *alarm_models.UsbModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New(my_errors.CurdClearFailed)
	}
	return nil
}
