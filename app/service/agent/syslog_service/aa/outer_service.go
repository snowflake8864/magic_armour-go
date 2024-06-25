package job_callback_syslog

import (
	"errors"
	"goskeleton/app/model/syslog/syslog_models"
)

func CreateOuterCurdFactory() *OuterCurd {
	return &OuterCurd{
		masterModel: syslog_models.CreateOuterFactory(),
	}
}

type OuterCurd struct {
	masterModel *syslog_models.OuterModel
}

// DeleteByIds 通过id集合删除日志
func (crd *OuterCurd) DeleteByIds(orgId int64, ids []int64) error {
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
func (crd *OuterCurd) ClearByAgents(orgId int64, agentIds []string) error {
	conditions := syslog_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("清空失败")
	}
	return nil
}
