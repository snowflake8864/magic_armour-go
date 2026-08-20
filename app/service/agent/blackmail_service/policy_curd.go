package blackmail_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_blackmail"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/blackmail_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
)

func CreatePolicyCurdFactory() *PolicyCurd {
	return &PolicyCurd{masterModel: blackmail_models.CreatePolicyFactory()}
}

type PolicyCurd struct {
	masterModel *blackmail_models.PolicyModel
}

// 新增任务
func (crd *PolicyCurd) addTask(orgId int64, agentId string, remark string) {
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypeBlackmailPolicyDownload, remark)
}

// Add 新增
func (crd *PolicyCurd) Add(param *blackmail_models.PolicyModel) error {
	param.Id = variable.SnowFlake.GetId()
	param.TypeVal = path_helper.ToStandard(param.TypeVal)
	param.OpenSwitch = enums.SwitchOff
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt
	if param.Type == enum_blackmail.TypeDir {
		param.GroupId = 0
	}

	//验证数据
	conditions := blackmail_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
		TypeVal: param.TypeVal,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New("已存在相同策略，不能重复添加")
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	return nil
}

// Edit 修改
func (crd *PolicyCurd) Edit(param *blackmail_models.PolicyModel, isRoot int) error {
	param.TypeVal = path_helper.ToStandard(param.TypeVal)
	param.UpdatedAt = date_helper.CurDatetime()
	if param.Type == enum_blackmail.TypeDir {
		param.GroupId = 0
	}

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(blackmail_models.SearchConditions{
		OrgId:   param.OrgId,
		Id:      param.Id,
		AgentId: param.AgentId,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	} else if hashModel.TypeVal != param.TypeVal {
		conditions := blackmail_models.SearchConditions{
			OrgId:   hashModel.OrgId,
			AgentId: hashModel.AgentId,
			TypeVal: param.TypeVal,
		}
		if total := crd.masterModel.Count(conditions); total > 0 {
			return errors.New("已存在相同策略，不能重复添加")
		}
	}

	//修改
	updates := make(map[string]interface{})
	updates["type"] = param.Type
	updates["type_val"] = param.TypeVal
	updates["group_id"] = param.GroupId
	updates["processes"] = param.Processes
	updates["remark"] = param.Remark
	updates["open_switch"] = enums.SwitchOff
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(blackmail_models.SearchConditions{
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
func (crd *PolicyCurd) SetSwitch(param *blackmail_models.PolicyModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(blackmail_models.SearchConditions{
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
	if isOk, _ := crd.masterModel.UpdateByConditions(blackmail_models.SearchConditions{
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
func (crd *PolicyCurd) DeleteById(param *blackmail_models.PolicyModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(blackmail_models.SearchConditions{
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
	if !crd.masterModel.DeleteByConditions(blackmail_models.SearchConditions{
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
	crd.masterModel.DeleteByConditions(blackmail_models.SearchConditions{
		OrgId: orgId,
	})
}
