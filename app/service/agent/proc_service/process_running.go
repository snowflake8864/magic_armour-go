package proc_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/task_service"
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

// ExtractDlls 提取模块
func (crd *ProcessRunning) ExtractDlls(orgId int64, agentId string, procHash string, procId int) error {

	hasModel := crd.masterModel.GetItemByOrgAndAgentAndHashAndProcId(orgId, agentId, procHash, procId)
	if hasModel == nil || hasModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	}

	updates := map[string]interface{}{
		"dll_status": enum_process.CollectDllStatusBegin,
		"dll_number": 0,
		"dll_list":   "",
	}
	isOk, rows := crd.masterModel.UpdateByOrgAndAgentAndHashAndProcId(orgId, agentId, procHash, procId, updates)
	if !isOk {
		return errors.New("修改数据错误")
	}

	if rows > 0 {

		task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypeProcessCollectDll, "")
	}
	return nil

}
