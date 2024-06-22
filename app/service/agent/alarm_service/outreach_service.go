package alarm_service

import (
	"errors"
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateOutreachCurdFactory() *OutreachCurd {
	return &OutreachCurd{
		masterModel: alarm_models.CreateOutreachFactory(),
	}
}

type OutreachCurd struct {
	alarmBase
	masterModel *alarm_models.OutreachModel
}

// Handle 处理
func (crd *OutreachCurd) Handle(param *alarm_models.OutreachModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		TargetIp: param.TargetIp,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New("保存失败")
	}

	return nil
}

// HandleAll 处理全部
func (crd *OutreachCurd) HandleAll(param *alarm_models.OutreachModel, agentIds []string) error {
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
func (crd *OutreachCurd) ClearByAgent(param *alarm_models.OutreachModel) error {
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
func (crd *OutreachCurd) ClearByAgents(param *alarm_models.OutreachModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("删除失败")
	}
	return nil
}
