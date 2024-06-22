package alarm_service

import (
	"errors"
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateTamperCurdFactory() *TamperCurd {
	return &TamperCurd{
		masterModel: alarm_models.CreateTamperFactory(),
	}
}

type TamperCurd struct {
	alarmBase
	masterModel *alarm_models.TamperModel
}

// Handle 处理
func (crd *TamperCurd) Handle(param *alarm_models.TamperModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New("保存失败")
	}

	return nil
}

// HandleAll 处理全部
func (crd *TamperCurd) HandleAll(param *alarm_models.TamperModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New("保存失败")
	}
	return nil
}

// ClearByAgent 清空终端下告警
func (crd *TamperCurd) ClearByAgent(param *alarm_models.TamperModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("清空失败")
	}
	return nil
}

// ClearByAgents 清空多终端下告警
func (crd *TamperCurd) ClearByAgents(param *alarm_models.TamperModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("清空失败")
	}
	return nil
}
