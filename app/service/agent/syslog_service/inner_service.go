package syslog_service

import (
	"errors"
	"goskeleton/app/model/syslog/syslog_models"
)

func CreateInnerCurdFactory() *InnerCurd {
	return &InnerCurd{
		masterModel: syslog_models.CreateInnerFactory(),
	}
}

type InnerCurd struct {
	masterModel *syslog_models.InnerModel
}

// DeleteByIds 通过id集合删除日志
func (crd *InnerCurd) DeleteByIds(orgId int64, ids []int64) error {
	conditions := syslog_models.SearchConditions{
		OrgId: orgId,
		Ids:   ids,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("删除失败")
	}
	return nil
}

// ClearByAgents 清空终端下日志
func (crd *InnerCurd) ClearByAgents(orgId int64, agentIds []string) error {
	conditions := syslog_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("清空失败")
	}
	return nil
}
