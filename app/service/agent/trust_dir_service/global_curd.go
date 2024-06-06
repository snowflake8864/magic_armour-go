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
	param.OpenSwitch = enums.SwitchOn
	param.CreatedAt = date_helper.CurDatetime()
	param.UpdatedAt = param.CreatedAt
	if param.DeptId < 0 {
		param.DeptId = 0
	}

	conditions := trust_models.SearchConditions{
		OrgId:        param.OrgId,
		PlatformType: param.PlatformType,
		DeptId:       param.DeptId,
		DirFull:      param.Dir,
	}
	if total := crd.masterModel.Count(conditions); total > 0 {
		return errors.New("'已存在相同策略，不能重复添加")
	}

	err, deptIds := crd.getDeptIds(param.OrgId, param.DeptId)
	if err != nil {
		return err
	}

	if !crd.masterModel.InsertItem(param) {
		return errors.New("保存失败")
	}

	//
	crd.addTask(param.OrgId, param.PlatformType, deptIds)

	return nil
}
