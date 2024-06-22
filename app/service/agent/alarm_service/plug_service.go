package alarm_service

import (
	"errors"
	"goskeleton/app/model/alarm/alarm_models"
)

func CreatPlugCurdFactory() *PlugCurd {
	return &PlugCurd{
		masterModel: alarm_models.CreatePlugFactory(),
	}
}

type PlugCurd struct {
	alarmBase
	masterModel *alarm_models.PlugModel
}

// Handle 处理
func (crd *PlugCurd) Handle(param *alarm_models.PlugModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId: param.OrgId,
		Ip:    param.Ip,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New("保存失败")
	}

	return nil
}

// HandleAll 处理全部
func (crd *PlugCurd) HandleAll(param *alarm_models.PlugModel, agentIds []string) error {
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
func (crd *PlugCurd) ClearByAgent(param *alarm_models.PlugModel) error {
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
func (crd *PlugCurd) ClearByAgents(param *alarm_models.PlugModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("删除失败")
	}
	return nil
}
