package plug_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreatePolicyDynamicCurdFactory() *PolicyDynamicCurd {
	return &PolicyDynamicCurd{masterModel: plug_models.CreateDynamicFactory()}
}

type PolicyDynamicCurd struct {
	masterModel *plug_models.DynamicModel
}

// 新增任务
func (crd *PolicyDynamicCurd) addTask(orgId int64, agentId string, remark string) {
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypePluggingDownload, remark)
}

// DisableById 解封
func (crd *PolicyDynamicCurd) DisableById(param *plug_models.DynamicModel) error {

	//验证数据
	OpenSwitch := enums.SwitchOn
	hashModel := crd.masterModel.GetItemByConditions(plug_models.SearchConditions{
		OrgId:      param.OrgId,
		Id:         param.Id,
		AgentId:    param.AgentId,
		OpenSwitch: &OpenSwitch,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	}

	//修改
	now := date_helper.CurDatetime()
	updates := make(map[string]interface{})
	updates["end_at"] = now
	updates["open_switch"] = enums.SwitchOff
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = now
	if isOk, _ := crd.masterModel.UpdateByConditions(plug_models.SearchConditions{
		OrgId:   hashModel.OrgId,
		Id:      hashModel.Id,
		AgentId: param.AgentId,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	//新增任务
	crd.addTask(hashModel.OrgId, hashModel.AgentId, "解封")

	return nil
}

// DisableByAgentAndIds 解封
func (crd *PolicyDynamicCurd) DisableByAgentAndIds(param *plug_models.DynamicModel, ids []int64) error {
	OpenSwitch := enums.SwitchOn
	now := date_helper.CurDatetime()
	//修改
	updates := make(map[string]interface{})
	updates["end_at"] = now
	updates["open_switch"] = enums.SwitchOff
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = now
	if isOk, _ := crd.masterModel.UpdateByConditions(plug_models.SearchConditions{
		OrgId:      param.OrgId,
		AgentId:    param.AgentId,
		Ids:        ids,
		OpenSwitch: &OpenSwitch,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	//新增任务
	crd.addTask(param.OrgId, param.AgentId, "终端批量解封")

	return nil
}

// DisableByAgentsAndIds 解封
func (crd *PolicyDynamicCurd) DisableByAgentsAndIds(param *plug_models.DynamicModel, ids []int64, agentIds []string) error {
	OpenSwitch := enums.SwitchOn
	now := date_helper.CurDatetime()
	//修改
	updates := make(map[string]interface{})
	updates["end_at"] = now
	updates["open_switch"] = enums.SwitchOff
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = now
	if isOk, _ := crd.masterModel.UpdateByConditions(plug_models.SearchConditions{
		OrgId:      param.OrgId,
		AgentIds:   agentIds,
		Ids:        ids,
		OpenSwitch: &OpenSwitch,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	//任务
	task_service.CreateCurdFactory().AddTaskForAgents(param.OrgId, agentIds, enum_task.TypePluggingDownload, "批量解封")

	return nil
}

// DeleteById 删除
func (crd *PolicyDynamicCurd) DeleteById(param *plug_models.DynamicModel) error {

	//验证数据
	OpenSwitch := enums.SwitchOff
	hashModel := crd.masterModel.GetItemByConditions(plug_models.SearchConditions{
		OrgId:      param.OrgId,
		Id:         param.Id,
		AgentId:    param.AgentId,
		OpenSwitch: &OpenSwitch,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	}

	//删除
	if !crd.masterModel.DeleteByConditions(plug_models.SearchConditions{
		OrgId:   hashModel.OrgId,
		Id:      hashModel.Id,
		AgentId: hashModel.AgentId,
	}) {
		return errors.New("删除失败")
	}

	return nil
}

// DeleteByAgentAndIds 删除
func (crd *PolicyDynamicCurd) DeleteByAgentAndIds(param *plug_models.DynamicModel, ids []int64) error {

	//删除
	if !crd.masterModel.DeleteByConditions(plug_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
		Ids:     ids,
	}) {
		return errors.New("删除失败")
	}

	return nil
}

// DeleteByAgentsAndIds 删除
func (crd *PolicyDynamicCurd) DeleteByAgentsAndIds(param *plug_models.DynamicModel, ids []int64, agentIds []string) error {

	conditions := plug_models.SearchConditions{
		OrgId:    param.OrgId,
		AgentIds: agentIds,
		Ids:      ids,
	}

	//删除
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("删除失败")
	}

	return nil
}
