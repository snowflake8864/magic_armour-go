package proc_service

import (
	"goskeleton/app/model/agent/proc_models"
)

func CreateDllRunningFactory() *DllRunning {
	return &DllRunning{
		masterModel: proc_models.CreateDllRunningFactory(),
	}
}

type DllRunning struct {
	ProcessBase
	masterModel *proc_models.DllRunningModel
}

// DeleteByAgents 通过终端集合删除
func (crd *DllRunning) DeleteByAgents(orgId int64, agentIds []string) bool {
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}
	return crd.masterModel.DeleteByConditions(condition)
}

// DeleteByAgent 通过终端删除
func (crd *DllRunning) DeleteByAgent(orgId int64, agentId string) bool {

	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	}
	return crd.masterModel.DeleteByConditions(condition)
}
