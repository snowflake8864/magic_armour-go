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

// DeleteByAgents 批量把终端下未设置变成白名单
func (crd *ProcessRunning) DeleteByAgents(orgId int64, agentIds []string) bool {
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}
	return crd.masterModel.DeleteByConditions(condition)
}
