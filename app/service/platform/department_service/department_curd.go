package department_service

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_platform"
)

func CreateCurdFactory() *DepartmentCurd {
	return &DepartmentCurd{
		masterModel: model_platform.CreateDepartmentFactory(),
	}
}

type DepartmentCurd struct {
	masterModel *model_platform.DepartmentModel
}

// AutoAllocate 自动分配
func (crd *DepartmentCurd) AutoAllocate(adminId int, agentId, agentIp string) {
	variable.ZapLog.Info("自动分配组织架构----------------")
}
