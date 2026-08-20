package plug_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/plug_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateGlobalIpCurdFactory() *GlobalIpCurd {
	return &GlobalIpCurd{masterModel: plug_models.CreateIpGlobalFactory()}
}

type GlobalIpCurd struct {
	masterModel *plug_models.IpGlobalModel
}

// 新增任务
func (crd *GlobalIpCurd) addTask(orgId int64, remark string) {
	isOnline := enums.IsYes
	switchOn := enums.SwitchOn
	conditions := agent_models.SearchConditions{
		OrgId:      orgId,
		IsOnline:   &isOnline,
		PlugSwitch: &switchOn,
	}
	if ids := agent_models.CreateAgentFactory().GetIdsByByConditions(conditions); ids != nil && len(ids) > 0 {
		task_service.CreateCurdFactory().AddTaskForAgents(orgId, ids, enum_task.TypePlugIpBlackDownload, remark)
	}
}

// Add 新增
func (crd *GlobalIpCurd) Add(param *plug_models.IpGlobalModel) error {
	param.Id = variable.SnowFlake.GetId()
	param.Direction = enum_plug.DirectionAll
	param.OpenSwitch = enums.SwitchOn
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt
	if param.EffectMode == enum_plug.EffectModeAlways {
		param.EffectStartAt = ""
		param.EffectEndAt = ""
	}

	//验证数据
	conditions := plug_models.SearchConditions{
		OrgId: param.OrgId,
		Ip:    param.Ip,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New(my_errors.CurdSamePolicyExistCannotAdd)
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	if param.Type == enum_plug.IpTypeBlack && param.OpenSwitch == enums.SwitchOn {
		//新增任务
		crd.addTask(param.OrgId, "Global Add")
	}

	return nil
}

// Edit 修改
func (crd *GlobalIpCurd) Edit(param *plug_models.IpGlobalModel, isRoot int) error {
	param.UpdatedAt = date_helper.CurDatetime()

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(plug_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New(my_errors.CurdOnlyModifyOneselfAdded)
	} else if hashModel.Ip != param.Ip {
		conditions := plug_models.SearchConditions{
			OrgId: hashModel.OrgId,
			Ip:    param.Ip,
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
	updates["effect_mode"] = param.EffectMode
	if param.EffectMode == enum_plug.EffectModeAlways {
		updates["effect_start_at"] = ""
		updates["effect_end_at"] = ""
	} else {
		updates["effect_start_at"] = param.EffectStartAt
		updates["effect_end_at"] = param.EffectEndAt
	}

	if isOk, _ := crd.masterModel.UpdateByConditions(plug_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	if hashModel.Type == enum_plug.IpTypeBlack && hashModel.OpenSwitch == enums.SwitchOn {
		//新增任务
		crd.addTask(hashModel.OrgId, "Global Edit")
	}

	return nil
}

// SetSwitch 开关
func (crd *GlobalIpCurd) SetSwitch(param *plug_models.IpGlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(plug_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
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
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//新增任务
	if hashModel.Type == enum_plug.IpTypeBlack {
		crd.addTask(hashModel.OrgId, "Global Switch")
	}

	return nil
}

// DeleteById 删除
func (crd *GlobalIpCurd) DeleteById(param *plug_models.IpGlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(plug_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New(my_errors.CurdOnlyDeleteOneselfAdded)
	}

	//删除
	if !crd.masterModel.DeleteByConditions(plug_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}) {
		return errors.New(my_errors.CurdDeleteFailed)
	}

	//新增任务
	if hashModel.Type == enum_plug.IpTypeBlack {
		crd.addTask(hashModel.OrgId, "Global Delete")
	}

	return nil
}

// DeleteByAgentAndType 删除
func (crd *GlobalIpCurd) DeleteByAgentAndType(param *plug_models.IpPolicyModel, isRoot int) error {
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
		crd.addTask(param.OrgId, "Global Clear")
	}

	return nil
}
