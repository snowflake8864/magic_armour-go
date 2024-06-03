package proc_service

import (
	"errors"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/task_service"
)

func CreateProcessCurdFactory() *ProcessCurd {
	return &ProcessCurd{
		masterModel: proc_models.CreateProcessFactory(),
	}
}

type ProcessCurd struct {
	ProcessBase
	masterModel *proc_models.ProcessModel
}

// BatchUnsetToWhiteByAgents 批量把终端下未设置变成白名单
func (crd *ProcessCurd) BatchUnsetToWhiteByAgents(orgId int64, agentIds []string) error {
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
		Status:   enum_process.ProcStatusUnset,
	}
	updates := map[string]interface{}{
		"self_status":    enum_process.RelyStatusWhite,
		"status":         enum_process.ProcStatusWhite,
		"sync_status_at": 0,
	}

	isOk, rows := crd.masterModel.UpdateByConditions(condition, updates)
	if !isOk {
		return errors.New("修改失败")
	}

	if rows > 0 {
		types := []int{
			enum_task.TypeProcessWhiteDownload,
			enum_task.TypeProcessBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgents(orgId, agentIds, types, "批量把终端下未设置变成白名单")
	}

	return nil
}
