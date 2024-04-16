package department_service

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_organization"
)

func CreateCurdFactory() *DepartmentCurd {
	return &DepartmentCurd{
		masterModel: model_organization.CreateDepartmentFactory(),
	}
}

type DepartmentCurd struct {
	masterModel *model_organization.DepartmentModel
}

// GetItemById 通过id获取一条记录
func (crd *DepartmentCurd) GetItemById(id int64) (model *model_organization.DepartmentModel) {
	model, _ = crd.masterModel.GetItemById(id)
	return model
}

// GetMapsByIds 通过ids获取多条记录, 有映射关系
func (crd *DepartmentCurd) GetMapsByIds(ids []int64) map[int64]*model_organization.DepartmentModel {
	models, err := crd.masterModel.GetItemsByIds(ids)
	if err != nil {
		return nil
	}
	if models == nil {
		return nil
	}
	result := make(map[int64]*model_organization.DepartmentModel)
	for _, model := range models {
		result[model.Id] = model
	}
	return result
}

// AutoAllocate 自动分配
func (crd *DepartmentCurd) AutoAllocate(orgId int64, agentId, agentIp string) {
	variable.ZapLog.Info("自动分配组织架构----------------")
}
