package plug_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/agent_models"
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
	updates := make(map[string]interface{})
	updates["end_at"] = date_helper.CurDatetime()
	updates["open_switch"] = enums.SwitchOff
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
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

// DisableByAgent 解封
func (crd *PolicyDynamicCurd) DisableByAgent(param *plug_models.DynamicModel) error {
	OpenSwitch := enums.SwitchOn
	//修改
	updates := make(map[string]interface{})
	updates["end_at"] = date_helper.CurDatetime()
	updates["open_switch"] = enums.SwitchOff
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(plug_models.SearchConditions{
		OrgId:      param.OrgId,
		AgentId:    param.AgentId,
		OpenSwitch: &OpenSwitch,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	//新增任务
	crd.addTask(param.OrgId, param.AgentId, "终端全部解封")

	return nil
}

// DisableByAgents 解封
func (crd *PolicyDynamicCurd) DisableByAgents(param *plug_models.DynamicModel, isRoot int, agentIds []string) error {
	OpenSwitch := enums.SwitchOn
	conditions := plug_models.SearchConditions{
		OrgId:      param.OrgId,
		OpenSwitch: &OpenSwitch,
	}
	if isRoot == enums.IsNot {
		if agentIds == nil || len(agentIds) == 0 {
			return nil
		}
		conditions.AgentIds = agentIds
	}

	//修改
	updates := make(map[string]interface{})
	updates["end_at"] = date_helper.CurDatetime()
	updates["open_switch"] = enums.SwitchOff
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	isOk, rows := crd.masterModel.UpdateByConditions(conditions, updates)
	if !isOk {
		return errors.New("保存失败")
	}

	//任务
	if rows > 0 {
		isOnline := enums.IsYes
		conditions1 := agent_models.SearchConditions{
			OrgId:      param.OrgId,
			IsOnline:   &isOnline,
			PlugSwitch: &OpenSwitch,
		}
		if isRoot == enums.IsNot {
			conditions.AgentIds = agentIds
		}
		if ids := agent_models.CreateAgentFactory().GetIdsByByConditions(conditions1); ids != nil || len(ids) > 0 {
			task_service.CreateCurdFactory().AddTaskForAgents(param.OrgId, ids, enum_task.TypePluggingDownload, "全部解封")
		}
	}

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
		OrgId:      hashModel.OrgId,
		Id:         hashModel.Id,
		AgentId:    hashModel.AgentId,
		OpenSwitch: &OpenSwitch,
	}) {
		return errors.New("删除失败")
	}

	return nil
}

// DeleteByAgent 删除
func (crd *PolicyDynamicCurd) DeleteByAgent(param *plug_models.DynamicModel) error {

	OpenSwitch := enums.SwitchOff
	//删除
	if !crd.masterModel.DeleteByConditions(plug_models.SearchConditions{
		OrgId:      param.OrgId,
		AgentId:    param.AgentId,
		OpenSwitch: &OpenSwitch,
	}) {
		return errors.New("删除失败")
	}

	return nil
}

// DeleteByAgents 删除
func (crd *PolicyDynamicCurd) DeleteByAgents(param *plug_models.DynamicModel, isRoot int, agentIds []string) error {

	OpenSwitch := enums.SwitchOff
	conditions := plug_models.SearchConditions{
		OrgId:      param.OrgId,
		OpenSwitch: &OpenSwitch,
	}
	if isRoot == enums.IsNot {
		if agentIds == nil || len(agentIds) == 0 {
			return nil
		}
		conditions.AgentIds = agentIds
	}
	//删除
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("删除失败")
	}

	return nil
}
