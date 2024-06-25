package syslog_service

import (
	"errors"
	"goskeleton/app/model/syslog/syslog_models"
)

func CreateProcessCurdFactory() *ProcessCurd {
	return &ProcessCurd{
		masterModel: syslog_models.CreateProcessFactory(),
	}
}

type ProcessCurd struct {
	masterModel *syslog_models.ProcessModel
}

// DeleteByIds 通过id集合删除日志
func (crd *ProcessCurd) DeleteByIds(orgId int64, ids []int64) error {
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
func (crd *ProcessCurd) ClearByAgents(orgId int64, agentIds []string) error {
	conditions := syslog_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("清空失败")
	}
	return nil
}
