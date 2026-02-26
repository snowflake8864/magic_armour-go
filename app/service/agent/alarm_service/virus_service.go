package alarm_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_alarm"
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateVirusCurdFactory() *VirusCurd {
	return &VirusCurd{
		masterModel: alarm_models.CreateVirusFactory(),
	}
}

type VirusCurd struct {
	alarmBase
	masterModel *alarm_models.VirusModel
}

// Handle 处理
func (crd *VirusCurd) Handle(param *alarm_models.VirusModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
		Id:      param.Id,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New("保存失败")
	}

	return nil
}

// HandleAll 处理全部
func (crd *VirusCurd) HandleAll(param *alarm_models.VirusModel, agentIds []string) error {
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
func (crd *VirusCurd) ClearByAgent(param *alarm_models.VirusModel) error {
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
func (crd *VirusCurd) ClearByAgents(param *alarm_models.VirusModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("清空失败")
	}
	return nil
}
