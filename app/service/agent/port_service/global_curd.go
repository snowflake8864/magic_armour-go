package port_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_port"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/port_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateGlobalCurdFactory() *GlobalCurd {
	return &GlobalCurd{masterModel: port_models.CreateVirtualGlobalFactory()}
}

type GlobalCurd struct {
	PortBase
	masterModel *port_models.VirtualGlobalModel
}

// 新增任务
func (crd *GlobalCurd) addTask(orgId int64, deptIds []int64, remark string) {
	isOnline := enums.IsYes
	switchOn := enums.SwitchOn
	conditions := agent_models.SearchConditions{
		OrgId:      orgId,
		DeptIds:    deptIds,
		IsOnline:   &isOnline,
		PortSwitch: &switchOn,
	}
	if ids := agent_models.CreateAgentFactory().GetIdsByByConditions(conditions); ids != nil && len(ids) > 0 {
		task_service.CreateCurdFactory().AddTaskForAgents(orgId, ids, enum_task.TypePortVirtualPolicyDownload, remark)
	}
}

// 获取组织架构以及子数据
func (crd *GlobalCurd) getDeptIds(orgId int64, deptId int64) (error, []int64) {
	var ids []int64
	if deptId > 0 {
		//获取组织架构
		ids = department_service.CreateCurdFactory().GetIdsSelfWithChildById(orgId, deptId)
		if ids == nil || len(ids) == 0 {
			return errors.New("组织架构不存在"), ids
		}
	}
	return nil, ids
}

// 获取已存在的虚开端口列表
func (crd *GlobalCurd) getExistsPorts(orgId int64, deptIds []int64, excludeId int64) []string {
	deptIds = append(deptIds, 0)
	conditions := port_models.SearchConditions{
		OrgId:   orgId,
		DeptIds: deptIds,
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
func (crd *GlobalCurd) Add(param *port_models.VirtualGlobalModel) error {
	param.Id = variable.SnowFlake.GetId()
	param.OpenSwitch = enums.SwitchOff
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt
	if param.DeptId < 0 {
		param.DeptId = 0
	}
	if param.Type == enum_port.IpTypeIpv4 {
		param.SourceIp = enum_port.Ipv4BehalfAll
	} else {
		param.SourceIp = enum_port.Ipv6BehalfAll
	}

	if err := crd.checkSourcePort(param.SourcePort); err != nil {
		return err
	}

	//组织架构
	err, deptIds := crd.getDeptIds(param.OrgId, param.DeptId)
	if err != nil {
		return err
	}
	if err = crd.checkPortUnique(param.SourcePort, crd.getExistsPorts(param.OrgId, deptIds, 0)); err != nil {
		return err
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}

	return nil
}

// Edit 修改
func (crd *GlobalCurd) Edit(param *port_models.VirtualGlobalModel, isRoot int) error {
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt
	if param.DeptId < 0 {
		param.DeptId = 0
	}
	if param.Type == enum_port.IpTypeIpv4 {
		param.SourceIp = enum_port.Ipv4BehalfAll
	} else {
		param.SourceIp = enum_port.Ipv6BehalfAll
	}

	if err := crd.checkSourcePort(param.SourcePort); err != nil {
		return err
	}

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(port_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	}

	//组织架构
	err, deptIds := crd.getDeptIds(hashModel.OrgId, param.DeptId)
	if err != nil {
		return err
	}

	if hashModel.SourcePort != param.SourcePort {
		if err = crd.checkPortUnique(param.SourcePort, crd.getExistsPorts(param.OrgId, deptIds, param.Id)); err != nil {
			return err
		}
	}

	//修改
	updates := make(map[string]interface{})
	updates["dept_id"] = param.DeptId
	updates["type"] = param.Type
	updates["protocol"] = param.Protocol
	updates["source_ip"] = param.SourceIp
	updates["source_port"] = param.SourcePort
	updates["forward_ip"] = param.ForwardIp
	updates["forward_port"] = param.ForwardPort
	updates["alarm_level"] = param.AlarmLevel
	updates["open_switch"] = enums.SwitchOff
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(port_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	if hashModel.OpenSwitch == enums.SwitchOn {
		//新增任务
		crd.addTask(hashModel.OrgId, deptIds, "全局：编辑")
	}

	return nil
}

// SetSwitch 开关
func (crd *GlobalCurd) SetSwitch(param *port_models.VirtualGlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(port_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
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
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	//组织架构
	err, deptIds := crd.getDeptIds(hashModel.OrgId, hashModel.DeptId)
	if err != nil {
		return nil
	}

	//新增任务
	crd.addTask(hashModel.OrgId, deptIds, "全局：开关")

	return nil
}

// DeleteById 删除
func (crd *GlobalCurd) DeleteById(param *port_models.VirtualGlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(port_models.SearchConditions{
		OrgId: param.OrgId,
		Id:    param.Id,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能删除自己添加的策略(管理员除外)")
	}

	//删除
	if !crd.masterModel.DeleteByConditions(port_models.SearchConditions{
		OrgId: hashModel.OrgId,
		Id:    hashModel.Id,
	}) {
		return errors.New("删除失败")
	}

	//组织架构
	err, deptIds := crd.getDeptIds(hashModel.OrgId, hashModel.DeptId)
	if err != nil {
		return nil
	}

	//新增任务
	crd.addTask(hashModel.OrgId, deptIds, "全局：删除")

	return nil
}

// DeleteByOrg 删除
func (crd *GlobalCurd) DeleteByOrg(orgId int64) {
	crd.masterModel.DeleteByConditions(port_models.SearchConditions{
		OrgId: orgId,
	})
}
