package proc_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/task_service"
)

func CreateDllCurdFactory() *DllCurd {
	return &DllCurd{
		masterModel: proc_models.CreateDllFactory(),
	}
}

type DllCurd struct {
	ProcessBase
	masterModel *proc_models.DllModel
}

// SetStatus 设置状态
func (crd *DllCurd) getOne(orgId int64, agentId string, hash string) *proc_models.DllModel {
	model := crd.masterModel.GetItemByOrgAndAgentAndHash(orgId, agentId, hash)
	if model == nil || model.OrgId == 0 {
		return nil
	} else {
		return model
	}
}

// BatchUnsetToWhiteByAgents 批量把终端下未设置变成白名单
func (crd *DllCurd) BatchUnsetToWhiteByAgents(orgId int64, agentIds []string) error {
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
			enum_task.TypeDllWhiteDownload,
			enum_task.TypeDllBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgents(orgId, agentIds, types, "批量把终端下未设置变成白名单")
	}

	return nil
}

// SetStatus 设置状态
func (crd *DllCurd) SetStatus(orgId int64, agentId string, hash string, status int) error {
	hasModel := crd.getOne(orgId, agentId, hash)
	if hasModel == nil {
		return errors.New("数据不存在")
	}
	if isOk := enum_process.BtnChangeEnable(hasModel.GlobalStatus); isOk == 0 {
		return errors.New("当前数据不允许修改状态")
	}
	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
		DllHash: hash,
	}
	updates := crd.getUpdatesForSetStatus(status)
	isOk, rows := crd.masterModel.UpdateByConditions(condition, updates)
	if !isOk {
		return errors.New("修改失败")
	}
	if rows > 0 {
		types := []int{
			enum_task.TypeDllWhiteDownload,
			enum_task.TypeDllBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgent(orgId, agentId, types, "设置状态")
	}

	return nil
}

// Cancel 取消
func (crd *DllCurd) Cancel(orgId int64, agentId string, hash string) error {
	hasModel := crd.getOne(orgId, agentId, hash)
	if hasModel == nil {
		return errors.New("数据不存在")
	}
	if isOk := enum_process.BtnCancelEnable(hasModel.SelfStatus); isOk == 0 {
		return errors.New("当前数据不允许取消")
	}
	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
		DllHash: hash,
	}
	updates := crd.getUpdatesForCancel()
	isOk, rows := crd.masterModel.UpdateByConditions(condition, updates)
	if !isOk {
		return errors.New("修改失败")
	}
	if rows > 0 {
		types := []int{
			enum_task.TypeDllWhiteDownload,
			enum_task.TypeDllBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgent(orgId, agentId, types, "取消")
	}

	return nil
}

// BatchWhiteByHashes 批量白名单
func (crd *DllCurd) BatchWhiteByHashes(orgId int64, agentId string, hashes []string) error {

	condition := proc_models.SearchConditions{
		OrgId:     orgId,
		AgentId:   agentId,
		DllHashes: hashes,
		StatusList: []int{
			enum_process.ProcStatusUnset,
			enum_process.ProcStatusUnknown,
			enum_process.ProcStatusBlack,
			enum_process.ProcStatusWhite,
		},
	}
	updates := crd.getUpdatesForBatchWhite()
	isOk, rows := crd.masterModel.UpdateByConditions(condition, updates)
	if !isOk {
		return errors.New("修改失败")
	}
	if rows > 0 {
		types := []int{
			enum_task.TypeDllWhiteDownload,
			enum_task.TypeDllBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgent(orgId, agentId, types, "批量白名单")
	}

	return nil
}

// DeleteByAgent 通过终端删除
func (crd *DllCurd) DeleteByAgent(orgId int64, agentId string, task bool) error {

	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	}
	isOk := crd.masterModel.DeleteByConditions(condition)
	if !isOk {
		return errors.New("删除失败")
	}

	if task {
		types := []int{
			enum_task.TypeDllWhiteDownload,
			enum_task.TypeDllBlackDownload,
		}
		task_service.CreateCurdFactory().AddTasksForAgent(orgId, agentId, types, "清空")
	}
	return nil
}

// DeleteByOrg 通过企业删除
func (crd *DllCurd) DeleteByOrg(orgId int64) error {

	condition := proc_models.SearchConditions{
		OrgId: orgId,
	}
	isOk := crd.masterModel.DeleteByConditions(condition)
	if !isOk {
		return errors.New("删除失败")
	}

	return nil
}

// GlobalCallbackAddByHash 全局名单回调： 新增一条
func (crd *DllCurd) GlobalCallbackAddByHash(orgId int64, globalType int, hash string) {
	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		DllHash: hash,
	}
	updates := crd.getUpdatesForAddGlobal(globalType)
	crd.masterModel.UpdateByConditions(condition, updates)
}

// GlobalCallbackAddByHashes 全局名单回调： 新增多条
func (crd *DllCurd) GlobalCallbackAddByHashes(orgId int64, globalType int, hashes []string) {
	condition := proc_models.SearchConditions{
		OrgId:     orgId,
		DllHashes: hashes,
	}
	updates := crd.getUpdatesForAddGlobal(globalType)
	crd.masterModel.UpdateByConditions(condition, updates)
}

// GlobalCallbackDeleteByHash 全局名单回调： 新增多条
func (crd *DllCurd) GlobalCallbackDeleteByHash(orgId int64, hash string) {
	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		DllHash: hash,
	}
	updates := crd.getUpdatesForDeleteGlobal()
	crd.masterModel.UpdateByConditions(condition, updates)
}

// GlobalCallbackDeleteAll 全局名单回调： 删除全部
func (crd *DllCurd) GlobalCallbackDeleteAll(orgId int64) {
	condition := proc_models.SearchConditions{
		OrgId: orgId,
	}
	updates := crd.getUpdatesForDeleteGlobal()
	crd.masterModel.UpdateByConditions(condition, updates)
}

// ThreatenCallbackByHash 威胁情报回调： 一条
func (crd *DllCurd) ThreatenCallbackByHash(hash string, level int) {
	condition := proc_models.SearchConditions{
		DllHash: hash,
	}
	if updates := crd.getUpdatesForThreatenCallback(level); updates != nil && len(updates) > 0 {
		crd.masterModel.UpdateByConditions(condition, updates)
	}
}

// ThreatenCallbackByHashes 威胁情报回调： 多条
func (crd *DllCurd) ThreatenCallbackByHashes(hashes []string, level int) {
	condition := proc_models.SearchConditions{
		DllHashes: hashes,
	}
	if updates := crd.getUpdatesForThreatenCallback(level); updates != nil && len(updates) > 0 {
		crd.masterModel.UpdateByConditions(condition, updates)
	}
}

// ThreatenRepairStatus 威胁情报修复状态
func (crd *DllCurd) ThreatenRepairStatus(orgId int64, threatenStatus int) {
	switchOn := enums.SwitchOn
	condition := proc_models.SearchConditions{
		OrgId:               orgId,
		Status:              enum_process.ProcStatusUnset,
		AutoWhiteSwitch:     &switchOn,
		GrantThreatenSwitch: &switchOn,
		ThreatenStatus:      &threatenStatus,
	}
	if updates := crd.getUpdatesForThreatenRepair(threatenStatus); updates != nil && len(updates) > 0 {
		crd.masterModel.UpdateByConditions(condition, updates)
	}
}

// AutoWhiteSwitchOnAgent 终端自动白名单
func (crd *DllCurd) AutoWhiteSwitchOnAgent(orgId int64, agentId string, openSwitch int) {
	condition := proc_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	}
	if updates := crd.getUpdatesForAutoWhiteSwitch(openSwitch); updates != nil && len(updates) > 0 {
		crd.masterModel.UpdateByConditions(condition, updates)
	}
}
