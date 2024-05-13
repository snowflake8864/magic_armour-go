package proc_service

import (
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
)

func CreateGlobalCurdFactory() *GlobalCurd {
	return &GlobalCurd{
		masterModel: model_agent.CreateProcGlobalFactory(),
	}
}

type GlobalCurd struct {
	masterModel *model_agent.ProcGlobalModel
}

// GetOne 获取一条
func (crd *GlobalCurd) GetOne(orgId int64, hash string) (model *model_agent.ProcGlobalModel) {
	model, err := crd.masterModel.GetItemByOrgAndId(orgId, hash)
	if err == nil || helpers.StringIsEmpty(model.Id) {
		return nil
	}
	return model
}

// Add 添加
func (crd *GlobalCurd) Add(param *model_agent.ProcGlobalModel) (errCode int, err error) {

	return errCode, err
}
