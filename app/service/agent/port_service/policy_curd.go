package port_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_port"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreatePolicyCurdFactory() *PolicyCurd {
	return &PolicyCurd{masterModel: port_models.CreateVirtualPolicyFactory()}
}

type PolicyCurd struct {
	PortBase
	masterModel *port_models.VirtualPolicyModel
}

// 新增任务
func (crd *PolicyCurd) addTask(orgId int64, agentId string, remark string) {
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypePortVirtualPolicyDownload, remark)
}

// 获取已存在的虚开端口列表
func (crd *PolicyCurd) getExistsPorts(orgId int64, agentId string, excludeId int64) []string {
	conditions := port_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	}

	models := crd.masterModel.GetListByConditions(conditions)
	if models == nil || len(models) == 0 {
		return nil
	}
	var portList []string
	for _, model := range models {
		if model.Id != excludeId {
			portList = append(portList, model.SourcePort)
		}
	}

	return portList
}

// Add 新增
func (crd *PolicyCurd) Add(param *port_models.VirtualPolicyModel) error {
	param.Id = variable.SnowFlake.GetId()
	param.OpenSwitch = enums.SwitchOff
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt
	if param.Type == enum_port.IpTypeIpv4 {
		param.SourceIp = enum_port.Ipv4BehalfAll
	} else {
		param.SourceIp = enum_port.Ipv6BehalfAll
	}

	if param.ForwardType == enum_port.ForwardTypeOriginal {
		param.ForwardPort = "0"
	}

	if err := crd.checkSourcePort(param.SourcePort); err != nil {
		return err
	}
	if err := crd.checkPortUnique(param.SourcePort, crd.getExistsPorts(param.OrgId, param.AgentId, 0)); err != nil {
		return err
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	return nil
}

// Edit 修改
func (crd *PolicyCurd) Edit(param *port_models.VirtualPolicyModel, isRoot int) error {
	param.UpdatedAt = date_helper.CurDatetime()
	if param.Type == enum_port.IpTypeIpv4 {
		param.SourceIp = enum_port.Ipv4BehalfAll
	} else {
		param.SourceIp = enum_port.Ipv6BehalfAll
	}
	if param.ForwardType == enum_port.ForwardTypeOriginal {
		param.ForwardPort = "0"
	}

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(port_models.SearchConditions{
		OrgId:   param.OrgId,
		Id:      param.Id,
		AgentId: param.AgentId,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	} else if hashModel.SourcePort != param.SourcePort {
		if err := crd.checkSourcePort(param.SourcePort); err != nil {
			return err
		}
		if err := crd.checkPortUnique(param.SourcePort, crd.getExistsPorts(param.OrgId, param.AgentId, param.Id)); err != nil {
			return err
		}
	}

	//修改
	updates := make(map[string]interface{})
	updates["type"] = param.Type
	updates["protocol"] = param.Protocol
	updates["source_ip"] = param.SourceIp
	updates["source_port"] = param.SourcePort
	updates["forward_ip"] = param.ForwardIp
	updates["forward_port"] = param.ForwardPort
	updates["forward_type"] = param.ForwardType
	updates["alarm_level"] = param.AlarmLevel
	updates["open_switch"] = enums.SwitchOff
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(port_models.SearchConditions{
		OrgId:   hashModel.OrgId,
		Id:      hashModel.Id,
		AgentId: hashModel.AgentId,
	}, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	if hashModel.OpenSwitch == enums.SwitchOn {
		//新增任务
		crd.addTask(hashModel.OrgId, hashModel.AgentId, "编辑")
	}

	return nil
}

// SetSwitch 开关
func (crd *PolicyCurd) SetSwitch(param *port_models.VirtualPolicyModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(port_models.SearchConditions{
		OrgId:   param.OrgId,
		Id:      param.Id,
		AgentId: param.AgentId,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	}

	//修改
	updates := make(map[string]interface{})
	updates["open_switch"] = param.OpenSwitch
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(port_models.SearchConditions{
		OrgId:   hashModel.OrgId,
		Id:      hashModel.Id,
		AgentId: param.AgentId,
	}, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//新增任务
	crd.addTask(hashModel.OrgId, hashModel.AgentId, "开关")

	return nil
}

// DeleteById 删除
func (crd *PolicyCurd) DeleteById(param *port_models.VirtualPolicyModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(port_models.SearchConditions{
		OrgId:   param.OrgId,
		Id:      param.Id,
		AgentId: param.AgentId,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能删除自己添加的策略(管理员除外)")
	}

	//删除
	if !crd.masterModel.DeleteByConditions(port_models.SearchConditions{
		OrgId:   hashModel.OrgId,
		Id:      hashModel.Id,
		AgentId: hashModel.AgentId,
	}) {
		return errors.New("删除失败")
	}

	//新增任务
	crd.addTask(hashModel.OrgId, hashModel.AgentId, "删除")

	return nil
}

// DeleteByOrg 删除
func (crd *PolicyCurd) DeleteByOrg(orgId int64) {
	crd.masterModel.DeleteByConditions(port_models.SearchConditions{
		OrgId: orgId,
	})
}
