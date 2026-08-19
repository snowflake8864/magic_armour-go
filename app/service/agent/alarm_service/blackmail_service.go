package alarm_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateBlackmailCurdFactory() *BlackmailCurd {
	return &BlackmailCurd{
		masterModel: alarm_models.CreateBlackmailFactory(),
	}
}

type BlackmailCurd struct {
	alarmBase
	masterModel *alarm_models.BlackmailModel
}

// Handle 处理
func (crd *BlackmailCurd) Handle(param *alarm_models.BlackmailModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:        param.OrgId,
		AgentId:      param.AgentId,
		Id:           param.Id,
		HandleStatus: enum_alarm.HandleStatusNot,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	return nil
}

// HandleAll 处理全部
func (crd *BlackmailCurd) HandleAll(param *alarm_models.BlackmailModel, agentIds []string) error {
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
func (crd *BlackmailCurd) ClearByAgent(param *alarm_models.BlackmailModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New(my_errors.CurdClearFailed)
	}
	return nil
}

// ClearByAgents 清空终端下告警
func (crd *BlackmailCurd) ClearByAgents(param *alarm_models.BlackmailModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New(my_errors.CurdClearFailed)
	}
	return nil
}
