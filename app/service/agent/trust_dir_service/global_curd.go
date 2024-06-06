package trust_dir_service

import (
	"errors"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/trust_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
)

func CreateGlobalCurdFactory() *GlobalCurd {
	return &GlobalCurd{masterModel: trust_models.CreateTrustGlobalFactory()}
}

type GlobalCurd struct {
	masterModel *trust_models.TrustGlobalModel
}

// 新增任务
func (crd *GlobalCurd) addTask(orgId int64, platformType int, deptIds []int64) {
	conditions := agent_models.SearchConditions{
		OrgId:        orgId,
		PlatformType: platformType,
		DeptIds:      deptIds,
	}
	AgentFact := agent_models.CreateAgentFactory()
	models := AgentFact.GetListByConditions(conditions, AgentFact.ColumnSimplify())
	if models != nil || len(models) > 0 {
		for _, model := range models {
			task_service.CreateCurdFactory().AddTaskForAgent(orgId, model.AgentId, enum_task.TypeTrustDir, "新增全局")
		}
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

// Add 新增
func (crd *GlobalCurd) Add(param *trust_models.TrustGlobalModel) error {
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
	conditions := trust_models.SearchConditions{
		OrgId:        param.OrgId,
		PlatformType: param.PlatformType,
		DeptId:       param.DeptId,
		DirFull:      param.Dir,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New("已存在相同策略，不能重复添加")
	}

	//组织架构
	err, _ := crd.getDeptIds(param.OrgId, param.DeptId)
	if err != nil {
		return err
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}

	return nil
}

// Edit 修改
func (crd *GlobalCurd) Edit(param *trust_models.TrustGlobalModel, isRoot int) error {
	param.Dir = path_helper.ToStandard(param.Dir)
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt
	if param.DeptId < 0 {
		param.DeptId = 0
	}

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(trust_models.SearchConditions{
		OrgId:        param.OrgId,
		Id:           param.Id,
		PlatformType: param.PlatformType,
	})
	if hashModel == nil || hashModel.OrgId == 0 {
		return errors.New("数据不存在")
	} else if isRoot == enums.IsNot && hashModel.CreateId != param.CreateId {
		return errors.New("只能修改自己添加的策略(管理员除外)")
	} else if hashModel.Dir != param.Dir {
		conditions := trust_models.SearchConditions{
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
	if isOk, _ := crd.masterModel.UpdateByConditions(trust_models.SearchConditions{
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
func (crd *GlobalCurd) SetSwitch(param *trust_models.TrustGlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(trust_models.SearchConditions{
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
	if isOk, _ := crd.masterModel.UpdateByConditions(trust_models.SearchConditions{
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
func (crd *GlobalCurd) DeleteById(param *trust_models.TrustGlobalModel, isRoot int) error {

	//验证数据
	hashModel := crd.masterModel.GetItemByConditions(trust_models.SearchConditions{
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
	if !crd.masterModel.DeleteByConditions(trust_models.SearchConditions{
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
func (crd *GlobalCurd) DeleteByOrg(orgId int64) {
	crd.masterModel.DeleteByConditions(trust_models.SearchConditions{
		OrgId: orgId,
	})
}
