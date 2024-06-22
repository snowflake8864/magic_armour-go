package alarm_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateSelfCurdFactory() *SelfCurd {
	return &SelfCurd{
		masterModel: alarm_models.CreateSelfFactory(),
	}
}

type SelfCurd struct {
	alarmBase
	masterModel *alarm_models.SelfModel
}

// Handle 处理
func (crd *SelfCurd) Handle(param *alarm_models.SelfModel) error {
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
func (crd *SelfCurd) HandleAll(param *alarm_models.SelfModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:        param.OrgId,
		AgentIds:     agentIds,
		HandleStatus: enum_alarm.HandleStatusNot,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New("保存失败")
	}
	return nil
}

// ClearByAgent 清空终端下告警
func (crd *SelfCurd) ClearByAgent(param *alarm_models.SelfModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:        param.OrgId,
		AgentId:      param.AgentId,
		HandleStatus: enum_alarm.HandleStatusNot,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("清空失败")
	}
	return nil
}

// ClearByAgents 清空多终端下告警
func (crd *SelfCurd) ClearByAgents(param *alarm_models.SelfModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("清空失败")
	}
	return nil
}
