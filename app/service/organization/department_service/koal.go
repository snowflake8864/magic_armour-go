package department_service

import (
	"goskeleton/app/model/organization/department_models"
)

func CreateKoalFactory() *Koal {
	return &Koal{deptModel: department_models.CreateDepartmentFactory()}
}

type Koal struct {
	deptModel *department_models.DepartmentModel
}

func (crd *Koal) Sync(deptParam *department_models.DepartmentModel) error {

	//是否存在
	if hasModel, _ := crd.deptModel.GetItemById(deptParam.OrgId, deptParam.Id); hasModel == nil || hasModel.Id == 0 {

	} else {

	}

	return nil
}

func (crd *Koal) Delete(ids []int64) error {

	return nil
}
