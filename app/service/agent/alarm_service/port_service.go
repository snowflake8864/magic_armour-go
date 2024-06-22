package alarm_service

import (
	"errors"
	"goskeleton/app/model/alarm/alarm_models"
)

func CreatePortCurdFactory() *PortCurd {
	return &PortCurd{
		masterModel: alarm_models.CreatePortFactory(),
	}
}

type PortCurd struct {
	alarmBase
	masterModel *alarm_models.PortModel
}

// Handle 处理
func (crd *PortCurd) Handle(param *alarm_models.PortModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AttackIp: param.AttackIp,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New("保存失败")
	}

	return nil
}

// HandleAll 处理全部
func (crd *PortCurd) HandleAll(param *alarm_models.PortModel, agentIds []string) error {
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
func (crd *PortCurd) ClearByAgent(param *alarm_models.PortModel) error {
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
func (crd *PortCurd) ClearByAgents(param *alarm_models.PortModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("删除失败")
	}
	return nil
}
