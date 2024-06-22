package alarm_service

import (
	"errors"
	"goskeleton/app/model/alarm/alarm_models"
)

func CreateLoginCurdFactory() *LoginCurd {
	return &LoginCurd{
		masterModel: alarm_models.CreateLoginFactory(),
	}
}

type LoginCurd struct {
	alarmBase
	masterModel *alarm_models.LoginModel
}

// Handle 处理
func (crd *LoginCurd) Handle(param *alarm_models.LoginModel) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		SourceIp: param.SourceIp,
	}
	updates := crd.getUpdateForHandle(param.HandleRemark, param.HandleUser)
	if !crd.masterModel.UpdateByConditions(conditions, updates) {
		return errors.New("保存失败")
	}

	return nil
}

// HandleAll 处理全部
func (crd *LoginCurd) HandleAll(param *alarm_models.LoginModel, agentIds []string) error {
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
func (crd *LoginCurd) ClearByAgent(param *alarm_models.LoginModel) error {
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
func (crd *LoginCurd) ClearByAgents(param *alarm_models.LoginModel, agentIds []string) error {
	conditions := alarm_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("删除失败")
	}
	return nil
}
