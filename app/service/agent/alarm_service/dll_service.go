package alarm_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateDllCurdFactory() *DllCurd {
	return &DllCurd{
		masterModel: alarm_models.CreateDllFactory(),
	}
}

type DllCurd struct {
	alarmBase
	masterModel *alarm_models.DllModel
}

// Handle 处理
func (crd *DllCurd) Handle(param *alarm_models.DllModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:        param.OrgId,
		AgentId:      param.AgentId,
		DllHash:      param.DllHash,
		HandleStatus: enum_alarm.HandleStatusNot,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	return nil
}

// HandleAll 处理全部
func (crd *DllCurd) HandleAll(param *alarm_models.DllModel, agentIds []string) error {
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
func (crd *DllCurd) ClearByAgent(param *alarm_models.DllModel) error {
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
func (crd *DllCurd) ClearByAgents(param *alarm_models.DllModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New(my_errors.CurdClearFailed)
	}
	return nil
}
