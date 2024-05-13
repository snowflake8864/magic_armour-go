package proc_service

import (
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
)

func CreateProcessCurdFactory() *ProcessCurd {
	return &ProcessCurd{
		masterModel: model_agent.CreateProcessFactory(),
	}
}

type ProcessCurd struct {
	ProcessBase
	masterModel *model_agent.ProcessModel
}

// GetOne 获取一条
func (crd *ProcessCurd) GetOne(orgId int64, agentId, hash string) (model *model_agent.ProcessModel) {
	model, err := crd.masterModel.GetItemByOrgAndAgentAndHash(orgId, agentId, hash)
	if err == nil || helpers.StringIsEmpty(model.ProcHash) {
		return nil
	}
	return model
}

// Add 添加
func (crd *ProcessCurd) Add(param *model_agent.ProcessModel) (errCode int, err error) {
	//如果存在， 忽略
	hasModel := crd.GetOne(param.OrgId, param.AgentId, param.ProcHash)
	if hasModel == nil || helpers.StringIsEmpty(hasModel.ProcHash) {
		return errCode, err
	}

	//威胁情报状态
	isExist := false
	if isExist, param.ThreatenStatus = crd.getThreatenStatus(param.ProcHash); !isExist {
		//加入威胁情报缓存
	}
	//全局名单状态
	param.GlobalStatus = crd.getGlobalStatus(param.OrgId, param.ProcHash)

	return errCode, err
}

func (crd *ProcessCurd) pushThreaten(param *model_agent.ProcessModel) {

}
