package proc_service

import (
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/utils/helpers"
)

func CreateProcessFactory() *ProcessCurd {
	return &ProcessCurd{
		masterModel: proc_models.CreateProcessFactory(),
	}
}

type ProcessCurd struct {
	ProcessBase
	masterModel *proc_models.ProcessModel
}

// GetOne 获取一条
func (crd *ProcessCurd) GetOne(orgId int64, agentId, hash string) (model *proc_models.ProcessModel) {
	model = crd.masterModel.GetItemByOrgAndAgentAndHash(orgId, agentId, hash)
	if helpers.StringIsEmpty(model.ProcHash) {
		return nil
	}
	return model
}
