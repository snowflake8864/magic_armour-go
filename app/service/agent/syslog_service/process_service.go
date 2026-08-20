package syslog_service

import (
	"errors"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/utils/helpers/date_helper"
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
func (crd *ProcessCurd) DeleteByIds(orgId int64, ids []int64, startTime, endTime string) error {
	conditions := syslog_models.SearchConditions{
		OrgId:       orgId,
		Ids:         ids,
		CreateAtGte: startTime,
		CreateAtLte: endTime,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New(my_errors.CurdDeleteFailed)
	}
	return nil
}

// ClearByAgents 清空终端下日志
func (crd *ProcessCurd) ClearByAgents(orgId int64, agentIds []string) error {
	conditions := syslog_models.SearchConditions{
		OrgId:       orgId,
		AgentIds:    agentIds,
		CreateAtLte: date_helper.CurDatetime(),
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New(my_errors.CurdClearFailed)
	}
	return nil
}
