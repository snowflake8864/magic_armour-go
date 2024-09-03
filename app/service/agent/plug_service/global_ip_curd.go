package plug_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_plug"
	"goskeleton/app/global/enums/enum_task"
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

	//验证数据
	conditions := plug_models.SearchConditions{
		OrgId: param.OrgId,
		Ip:    param.Ip,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New("已存在相同策略，不能重复添加")
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}

	if param.Type == enum_plug.IpTypeBlack && param.OpenSwitch == enums.SwitchOn {
		//新增任务
		crd.addTask(param.OrgId, "添加")
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
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	} else if hashModel.Ip != param.Ip {
		conditions := plug_models.SearchConditions{
			OrgId: hashModel.OrgId,
			Ip:    param.Ip,
		}
		if total := crd.masterModel.Count(conditions); total > 0 {
			return errors.New("已存在相同策略，不能重复添加")
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
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	if hashModel.Type == enum_plug.IpTypeBlack && hashModel.OpenSwitch == enums.SwitchOn {
		//新增任务
		crd.addTask(hashModel.OrgId, "编辑")
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
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
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
		return errors.New("保存失败")
	}

	//新增任务
	if hashModel.Type == enum_plug.IpTypeBlack {
		crd.addTask(hashModel.OrgId, "开关")
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
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New("只能删除自己添加的策略(管理员除外)")
	}

	//删除
	if !crd.masterModel.DeleteByConditions(plug_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}) {
		return errors.New("删除失败")
	}

	//新增任务
	if hashModel.Type == enum_plug.IpTypeBlack {
		crd.addTask(hashModel.OrgId, "删除")
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
		return errors.New("删除失败")
	}

	//新增任务
	if param.Type == enum_plug.IpTypeBlack {
		crd.addTask(param.OrgId, "删除全部")
	}

	return nil
}
