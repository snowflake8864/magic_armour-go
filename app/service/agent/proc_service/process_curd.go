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

// SetStatus 设置状态
func (crd *ProcessCurd) getOne(orgId int64, agentId string, hash string) *proc_models.ProcessModel {
	model := crd.masterModel.GetItemByOrgAndAgentAndHash(orgId, agentId, hash)
	if model == nil || model.OrgId == 0 {
		return nil
	} else {
		return model
	}
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

// SetStatus 设置状态
func (crd *ProcessCurd) SetStatus(orgId int64, agentId string, hash string, status int) error {
	hasModel := crd.getOne(orgId, agentId, hash)
	if hasModel == nil {
		return errors.New("数据不存在")
	}
	if isOk := enum_process.BtnChangeEnable(hasModel.GlobalStatus); isOk == 0 {
		return errors.New("当前数据不允许修改状态")
	}
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		AgentId:  agentId,
		ProcHash: hash,
	}
	updates := crd.getUpdatesForSetStatus(status)
	isOk, rows := crd.masterModel.UpdateByConditions(condition, updates)
	if !isOk {
		return errors.New("修改失败")
	}
	if rows > 0 {
		types := []int{
			enum_task.TypeProcessWhiteDownload,
			enum_task.TypeProcessBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgent(orgId, agentId, types, "设置状态")
	}

	return nil
}

// Cancel 取消
func (crd *ProcessCurd) Cancel(orgId int64, agentId string, hash string) error {
	hasModel := crd.getOne(orgId, agentId, hash)
	if hasModel == nil {
		return errors.New("数据不存在")
	}
	if isOk := enum_process.BtnCancelEnable(hasModel.SelfStatus); isOk == 0 {
		return errors.New("当前数据不允许取消")
	}
	condition := proc_models.SearchConditions{
		OrgId:    orgId,
		AgentId:  agentId,
		ProcHash: hash,
	}
	updates := crd.getUpdatesForCancel()
	isOk, rows := crd.masterModel.UpdateByConditions(condition, updates)
	if !isOk {
		return errors.New("修改失败")
	}
	if rows > 0 {
		types := []int{
			enum_task.TypeProcessWhiteDownload,
			enum_task.TypeProcessBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgent(orgId, agentId, types, "取消")
	}

	return nil
}

// ClearAll 清空全部
func (crd *ProcessCurd) ClearAll(orgId int64, agentId string) error {

	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	}
	isOk := crd.masterModel.DeleteByConditions(condition)
	if !isOk {
		return errors.New("删除失败")
	}

	types := []int{
		enum_task.TypeProcessWhiteDownload,
		enum_task.TypeProcessBlackDownload,
	}
	task_service.CreateCurdFactory().AddTasksForAgent(orgId, agentId, types, "清空")

	return nil
}
