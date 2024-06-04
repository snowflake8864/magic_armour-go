package proc_service

import (
	"errors"
	"goskeleton/app/model/agent/proc_models"
)

func CreateProcessRunningFactory() *ProcessRunning {
	return &ProcessRunning{
		masterModel: proc_models.CreateProcessRunningFactory(),
	}
}

type ProcessRunning struct {
	ProcessBase
	masterModel *proc_models.ProcessRunningModel
}

// DeleteByAgents 批量把终端下未设置变成白名单
func (crd *ProcessRunning) DeleteByAgents(orgId int64, agentIds []string) bool {
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}
	return crd.masterModel.DeleteByConditions(condition)
}

// ClearRunning 清空运行中
func (crd *ProcessRunning) ClearRunning(orgId int64, agentId string) error {

	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	}
	isOk := crd.masterModel.DeleteByConditions(condition)
	if !isOk {
		return errors.New("删除失败")
	}

	return nil
}
