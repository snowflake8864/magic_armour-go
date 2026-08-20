package plug_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreatePolicyIpCurdFactory() *PolicyIpCurd {
	return &PolicyIpCurd{masterModel: plug_models.CreateIpPolicyFactory()}
}

type PolicyIpCurd struct {
	masterModel *plug_models.IpPolicyModel
}

// 新增任务
func (crd *PolicyIpCurd) addTask(orgId int64, agentId string, remark string) {
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypePlugIpBlackDownload, remark)
}

// Add 新增
func (crd *PolicyIpCurd) Add(param *plug_models.IpPolicyModel) error {
	param.Id = variable.SnowFlake.GetId()
	param.Direction = enum_plug.DirectionAll
	param.OpenSwitch = enums.SwitchOn
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt

	//验证数据
	conditions := plug_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
		Ip:      param.Ip,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New(my_errors.CurdSamePolicyExistCannotAdd)
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	if param.Type == enum_plug.IpTypeBlack && param.OpenSwitch == enums.SwitchOn {
		//新增任务
		crd.addTask(param.OrgId, param.AgentId, "Add")
	}

	return nil
}

// Edit 修改
func (crd *PolicyIpCurd) Edit(param *plug_models.IpPolicyModel, isRoot int) error {
	param.UpdatedAt = date_helper.CurDatetime()

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(plug_models.SearchConditions{
		OrgId:   param.OrgId,
		Id:      param.Id,
		AgentId: param.AgentId,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New(my_errors.CurdOnlyModifyOneselfAdded)
	} else if hashModel.Ip != param.Ip {
		conditions := plug_models.SearchConditions{
			OrgId:   hashModel.OrgId,
			AgentId: hashModel.AgentId,
			Ip:      param.Ip,
		}
		if total := crd.masterModel.Count(conditions); total > 0 {
			return errors.New(my_errors.CurdSamePolicyExistCannotAdd)
		}
	}

	//修改
	updates := make(map[string]interface{})
	updates["ip"] = param.Ip
	updates["remark"] = param.Remark
	//updates["open_switch"] = enums.SwitchOff
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(plug_models.SearchConditions{
		OrgId:   hashModel.OrgId,
		Id:      hashModel.Id,
		AgentId: hashModel.AgentId,
	}, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	if hashModel.Type == enum_plug.IpTypeBlack && hashModel.OpenSwitch == enums.SwitchOn {
		//新增任务
		crd.addTask(hashModel.OrgId, hashModel.AgentId, "Edit")
	}

	return nil
}

// SetSwitch 开关
func (crd *PolicyIpCurd) SetSwitch(param *plug_models.IpPolicyModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(plug_models.SearchConditions{
		OrgId:   param.OrgId,
		Id:      param.Id,
		AgentId: param.AgentId,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New(my_errors.CurdOnlyModifyOneselfAdded)
	}

	//修改
	updates := make(map[string]interface{})
	updates["open_switch"] = param.OpenSwitch
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(plug_models.SearchConditions{
		OrgId:   hashModel.OrgId,
		Id:      hashModel.Id,
		AgentId: param.AgentId,
	}, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//新增任务
	if hashModel.Type == enum_plug.IpTypeBlack {
		crd.addTask(hashModel.OrgId, hashModel.AgentId, "Switch")
	}

	return nil
}

// DeleteById 删除
func (crd *PolicyIpCurd) DeleteById(param *plug_models.IpPolicyModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(plug_models.SearchConditions{
		OrgId:   param.OrgId,
		Id:      param.Id,
		AgentId: param.AgentId,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New(my_errors.CurdOnlyDeleteOneselfAdded)
	}

	//删除
	if !crd.masterModel.DeleteByConditions(plug_models.SearchConditions{
		OrgId:   hashModel.OrgId,
		Id:      hashModel.Id,
		AgentId: hashModel.AgentId,
	}) {
		return errors.New(my_errors.CurdDeleteFailed)
	}

	//新增任务
	if hashModel.Type == enum_plug.IpTypeBlack {
		crd.addTask(hashModel.OrgId, hashModel.AgentId, "Delete")
	}

	return nil
}

// DeleteByAgentAndType 删除
func (crd *PolicyIpCurd) DeleteByAgentAndType(param *plug_models.IpPolicyModel, isRoot int) error {
	conditions := plug_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
		IpType:  param.Type,
	}
	//不是管理员只能删除自己创建的
	if isRoot == enums.IsNot {
		conditions.CreateId = param.UpdateId
	}

	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New(my_errors.CurdDeleteFailed)
	}

	//新增任务
	if param.Type == enum_plug.IpTypeBlack {
		crd.addTask(param.OrgId, param.AgentId, "Clear")
	}

	return nil
}
