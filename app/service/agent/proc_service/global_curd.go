package proc_service

import (
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/utils/helpers/string_helper"
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
	model = crd.masterModel.GetItemByOrgAndId(orgId, hash)
	if model == nil || string_helper.IsEmpty(model.Id) {
		return nil
	}
	return model
}

// Add 添加
func (crd *GlobalCurd) Add(param *proc_models.ProcGlobalModel) (errCode int, err error) {

	return errCode, err
}
