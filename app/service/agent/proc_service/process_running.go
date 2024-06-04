package proc_service

import (
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

// DeleteByAgents 通过终端集合删除
func (crd *ProcessRunning) DeleteByAgents(orgId int64, agentIds []string) bool {
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}
	return crd.masterModel.DeleteByConditions(condition)
}

// DeleteByAgent 通过终端删除
func (crd *ProcessRunning) DeleteByAgent(orgId int64, agentId string) bool {

	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	}
	return crd.masterModel.DeleteByConditions(condition)
}
