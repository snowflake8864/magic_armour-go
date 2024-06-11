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
	return &PolicyCurd{masterModel: tamper_models.CreateTamperPolicyFactory()}
}

type PolicyCurd struct {
	masterModel *tamper_models.TamperPolicyModel
}

// 新增任务
func (crd *PolicyCurd) addTask(orgId int64, agentId string, remark string) {
	task_service.CreateCurdFactory().AddTaskForAgent(orgId, agentId, enum_task.TypeTamperPolicyDownload, remark)
}

// Add 新增
func (crd *PolicyCurd) Add(param *tamper_models.TamperPolicyModel) error {
	param.Id = variable.SnowFlake.GetId()
	param.Dir = path_helper.ToStandard(param.Dir)
	param.OpenSwitch = enums.SwitchOn
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

	//新增任务
	crd.addTask(param.OrgId, param.AgentId, "添加")

	return nil
}

// Edit 修改
func (crd *PolicyCurd) Edit(param *tamper_models.TamperPolicyModel, isRoot int) error {
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt
	if param.DeptId < 0 {
		param.DeptId = 0
	}

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(tamper_models.SearchConditions{
		OrgId:        param.OrgId,
		Id:           param.Id,
		PlatformType: param.PlatformType,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	} else if hashModel.Dir != param.Dir {
		conditions := tamper_models.SearchConditions{
			OrgId:        hashModel.OrgId,
			PlatformType: hashModel.PlatformType,
			DeptId:       param.DeptId,
			DirFull:      param.Dir,
		}
		if total := crd.masterModel.Count(conditions); total > 0 {
			return errors.New("已存在相同策略，不能重复添加")
		}
	}

	//组织架构
	err, deptIds := crd.getDeptIds(hashModel.OrgId, param.DeptId)
	if err != nil {
		return err
	}

	//修改
	updates := make(map[string]interface{})
	updates["name"] = param.Name
	updates["dir"] = param.Dir
	updates["is_extend"] = param.IsExtend
	updates["dept_id"] = param.DeptId
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(tamper_models.SearchConditions{
		OrgId:        hashModel.OrgId,
		Id:           hashModel.Id,
		PlatformType: hashModel.PlatformType,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	//新增任务
	crd.addTask(hashModel.OrgId, hashModel.PlatformType, deptIds)

	return nil
}

// SetSwitch 开关
func (crd *PolicyCurd) SetSwitch(param *tamper_models.TamperPolicyModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(tamper_models.SearchConditions{
		OrgId:        param.OrgId,
		Id:           param.Id,
		PlatformType: param.PlatformType,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	}

	//修改
	updates := make(map[string]interface{})
	updates["open_switch"] = param.OpenSwitch
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(tamper_models.SearchConditions{
		OrgId:        hashModel.OrgId,
		Id:           hashModel.Id,
		PlatformType: hashModel.PlatformType,
	}, updates); !isOk {
		return errors.New("保存失败")
	}

	//组织架构
	err, deptIds := crd.getDeptIds(hashModel.OrgId, hashModel.DeptId)
	if err != nil {
		return err
	}

	//新增任务
	crd.addTask(hashModel.OrgId, hashModel.PlatformType, deptIds)

	return nil
}

// DeleteById 删除
func (crd *PolicyCurd) DeleteById(param *tamper_models.TamperPolicyModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(tamper_models.SearchConditions{
		OrgId:        param.OrgId,
		Id:           param.Id,
		PlatformType: param.PlatformType,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能删除自己添加的策略(管理员除外)")
	}

	//修改
	if !crd.masterModel.DeleteByConditions(tamper_models.SearchConditions{
		OrgId:        hashModel.OrgId,
		Id:           hashModel.Id,
		PlatformType: hashModel.PlatformType,
	}) {
		return errors.New("删除失败")
	}

	//组织架构
	err, deptIds := crd.getDeptIds(hashModel.OrgId, hashModel.DeptId)
	if err != nil {
		return err
	}

	//新增任务
	crd.addTask(hashModel.OrgId, hashModel.PlatformType, deptIds)

	return nil
}

// DeleteByOrg 删除
func (crd *PolicyCurd) DeleteByOrg(orgId int64) {
	crd.masterModel.DeleteByConditions(tamper_models.SearchConditions{
		OrgId: orgId,
	})
}
