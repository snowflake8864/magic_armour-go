package tamper_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
)

func CreatePolicyCurdFactory() *PolicyCurd {
	return &PolicyCurd{masterModel: tamper_models.CreatePolicyFactory()}
}

type PolicyCurd struct {
	masterModel *tamper_models.PolicyModel
}

// 新增任务
func (crd *PolicyCurd) addTask(orgId int64, agentId string, remark string) {
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypeTamperPolicyDownload, remark)
}

// Add 新增
func (crd *PolicyCurd) Add(param *tamper_models.PolicyModel) error {
	param.Id = variable.SnowFlake.GetId()
	param.Dir = path_helper.ToStandard(param.Dir)
	param.OpenSwitch = enums.SwitchOff
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt

	//验证数据
	conditions := tamper_models.SearchConditions{
		OrgId:   param.OrgId,
		AgentId: param.AgentId,
		DirFull: param.Dir,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New("已存在相同策略，不能重复添加")
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}

	return nil
}

// Edit 修改
func (crd *PolicyCurd) Edit(param *tamper_models.PolicyModel, isRoot int) error {
	param.Dir = path_helper.ToStandard(param.Dir)
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(tamper_models.SearchConditions{
		OrgId:   param.OrgId,
		Id:      param.Id,
		AgentId: param.AgentId,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	} else if hashModel.Dir != param.Dir {
		conditions := tamper_models.SearchConditions{
			OrgId:   hashModel.OrgId,
			AgentId: hashModel.AgentId,
			DirFull: param.Dir,
		}
		if total := crd.masterModel.Count(conditions); total > 0 {
			return errors.New("已存在相同策略，不能重复添加")
		}
	}

	//修改
	updates := make(map[string]interface{})
	updates["dir"] = param.Dir
	updates["protect_rw"] = param.ProtectRw
	updates["include_suffix"] = param.IncludeSuffix
	updates["exclude_suffix"] = param.ExcludeSuffix
	updates["exclude_folder"] = param.ExcludeFolder
	updates["is_extend"] = param.IsExtend
	updates["processes"] = param.Processes
	updates["remark"] = param.Remark
	updates["open_switch"] = param.OpenSwitch
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(tamper_models.SearchConditions{
		OrgId:   hashModel.OrgId,
		Id:      hashModel.Id,
		AgentId: hashModel.AgentId,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	if hashModel.OpenSwitch == enums.SwitchOn {
		//新增任务
		crd.addTask(hashModel.OrgId, hashModel.AgentId, "编辑")
	}

	return nil
}

// SetSwitch 开关
func (crd *PolicyCurd) SetSwitch(param *tamper_models.PolicyModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(tamper_models.SearchConditions{
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
	if isOk, _ := crd.masterModel.UpdateByConditions(tamper_models.SearchConditions{
		OrgId:   hashModel.OrgId,
		Id:      hashModel.Id,
		AgentId: param.AgentId,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	//新增任务
	crd.addTask(hashModel.OrgId, hashModel.AgentId, "开关")

	return nil
}

// DeleteById 删除
func (crd *PolicyCurd) DeleteById(param *tamper_models.PolicyModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(tamper_models.SearchConditions{
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
	if !crd.masterModel.DeleteByConditions(tamper_models.SearchConditions{
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
	crd.masterModel.DeleteByConditions(tamper_models.SearchConditions{
		OrgId: orgId,
	})
}
