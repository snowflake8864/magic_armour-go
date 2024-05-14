package proc_service

import (
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/utils/helpers"
)

func CreateGlobalCurdFactory() *GlobalCurd {
	return &GlobalCurd{
		masterModel: proc_models.CreateProcGlobalFactory(),
	}
}

type GlobalCurd struct {
	masterModel *proc_models.ProcGlobalModel
}

// GetOne 获取一条
func (crd *GlobalCurd) GetOne(orgId int64, hash string) (model *proc_models.ProcGlobalModel) {
	model, err := crd.masterModel.GetItemByOrgAndId(orgId, hash)
	if err == nil || helpers.StringIsEmpty(model.Id) {
		return nil
	}
	return model
}

// Add 添加
func (crd *GlobalCurd) Add(param *proc_models.ProcGlobalModel) (errCode int, err error) {

	return errCode, err
}
