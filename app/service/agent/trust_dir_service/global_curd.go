package trust_dir_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/trust_dir_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
)

func CreateGlobalCurdFactory() *GlobalCurd {
	return &GlobalCurd{masterModel: trust_dir_models.CreateGlobalFactory()}
}

type GlobalCurd struct {
	masterModel *trust_dir_models.GlobalModel
}

// 新增任务
func (crd *GlobalCurd) addTask(orgId int64, platformType int, deptIds []int64, remark string) {
	isOnline := enums.IsYes
	conditions := agent_models.SearchConditions{
		OrgId:        orgId,
		PlatformType: platformType,
		DeptIds:      deptIds,
		IsOnline:     &isOnline,
	}
	if ids := agent_models.CreateAgentFactory().GetIdsByByConditions(conditions); ids != nil && len(ids) > 0 {
		task_service.CreateCurdFactory().AddTaskForAgents(orgId, ids, enum_task.TypeTrustDir, remark)
	}
}

// 获取组织架构以及子数据
func (crd *GlobalCurd) getDeptIds(orgId int64, deptId int64) (error, []int64) {
	var ids []int64
	if deptId > 0 {
		//获取组织架构
		ids = department_service.CreateCurdFactory().GetIdsSelfWithChildById(orgId, deptId)
		if ids == nil || len(ids) == 0 {
			return errors.New(my_errors.CurdDepartmentNotExist), ids
		}
	}
	return nil, ids
}

// Add 新增
func (crd *GlobalCurd) Add(param *trust_dir_models.GlobalModel) error {
	param.Id = variable.SnowFlake.GetId()
	param.Dir = path_helper.ToStandard(param.Dir)
	param.Type = enum_tamper.TypeFolder
	param.OpenSwitch = enums.SwitchOff
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt
	if param.DeptId < 0 {
		param.DeptId = 0
	}

	//验证数据
	conditions := trust_dir_models.SearchConditions{
		OrgId:        param.OrgId,
		PlatformType: param.PlatformType,
		DeptId:       param.DeptId,
		DirFull:      param.Dir,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New(my_errors.CurdSamePolicyExistCannotAdd)
	}

	//组织架构
	err, _ := crd.getDeptIds(param.OrgId, param.DeptId)
	if err != nil {
		return err
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	return nil
}

// Edit 修改
func (crd *GlobalCurd) Edit(param *trust_dir_models.GlobalModel, isRoot int) error {
	param.Dir = path_helper.ToStandard(param.Dir)
	param.UpdatedAt = date_helper.CurDatetime()
	if param.DeptId < 0 {
		param.DeptId = 0
	}

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(trust_dir_models.SearchConditions{
		OrgId:        param.OrgId,
		Id:           param.Id,
		PlatformType: param.PlatformType,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New(my_errors.CurdOnlyModifyOneselfAdded)
	} else if hashModel.Dir != param.Dir {
		conditions := trust_dir_models.SearchConditions{
			OrgId:        hashModel.OrgId,
			PlatformType: hashModel.PlatformType,
			DeptId:       param.DeptId,
			DirFull:      param.Dir,
		}
		if total := crd.masterModel.Count(conditions); total > 0 {
			return errors.New(my_errors.CurdSamePolicyExistCannotAdd)
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
	updates["open_switch"] = enums.SwitchOff
	updates["dept_id"] = param.DeptId
	updates["update_id"] = param.UpdateId
	updates["update_name"] = param.UpdateName
	updates["updated_at"] = date_helper.CurDatetime()
	if isOk, _ := crd.masterModel.UpdateByConditions(trust_dir_models.SearchConditions{
		OrgId:        hashModel.OrgId,
		Id:           hashModel.Id,
		PlatformType: hashModel.PlatformType,
	}, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	if hashModel.OpenSwitch == enums.SwitchOn {
		//新增任务
		crd.addTask(hashModel.OrgId, hashModel.PlatformType, deptIds, "全局:修改")
	}
	return nil
}

// SetSwitch 开关
func (crd *GlobalCurd) SetSwitch(param *trust_dir_models.GlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(trust_dir_models.SearchConditions{
		OrgId:        param.OrgId,
		Id:           param.Id,
		PlatformType: param.PlatformType,
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
	if isOk, _ := crd.masterModel.UpdateByConditions(trust_dir_models.SearchConditions{
		OrgId:        hashModel.OrgId,
		Id:           hashModel.Id,
		PlatformType: hashModel.PlatformType,
	}, updates); !isOk {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//组织架构
	err, deptIds := crd.getDeptIds(hashModel.OrgId, hashModel.DeptId)
	if err != nil {
		return nil
	}

	//新增任务
	crd.addTask(hashModel.OrgId, hashModel.PlatformType, deptIds, "全局:开关")

	return nil
}

// DeleteById 删除
func (crd *GlobalCurd) DeleteById(param *trust_dir_models.GlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(trust_dir_models.SearchConditions{
		OrgId:        param.OrgId,
		Id:           param.Id,
		PlatformType: param.PlatformType,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New(my_errors.CurdDataNotExist)
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.UpdateId {
		return errors.New(my_errors.CurdOnlyDeleteOneselfAdded)
	}

	//修改
	if !crd.masterModel.DeleteByConditions(trust_dir_models.SearchConditions{
		OrgId:        hashModel.OrgId,
		Id:           hashModel.Id,
		PlatformType: hashModel.PlatformType,
	}) {
		return errors.New(my_errors.CurdDeleteFailed)
	}

	//组织架构
	err, deptIds := crd.getDeptIds(hashModel.OrgId, hashModel.DeptId)
	if err != nil {
		return nil
	}

	//新增任务
	crd.addTask(hashModel.OrgId, hashModel.PlatformType, deptIds, "全局:删除")

	return nil
}
