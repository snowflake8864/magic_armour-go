package port_service

import (
	"goskeleton/app/model/agent/port_models"
)

func CreateWorkingFactory() *WorkingCurd {
	return &WorkingCurd{port_models.CreateWorkingFactory()}
}

type WorkingCurd struct {
	masterModel *port_models.WorkingModel
}

// BatchAdd 新增
func (crd *WorkingCurd) BatchAdd(orgId int64, agentId string, params *[]port_models.WorkingModel) bool {
	//先删除
	crd.DelByOrgAndAgent(orgId, agentId)

	return crd.masterModel.InsertBatch(params)
}

// DelByOrg 删除：企业
func (crd *WorkingCurd) DelByOrg(orgId int64) bool {
	search := port_models.SearchConditions{OrgId: orgId}
	return crd.masterModel.DeleteByConditions(search)
}

// DelByOrgAndAgent 删除：企业和终端
func (crd *WorkingCurd) DelByOrgAndAgent(orgId int64, agentId string) bool {
	search := port_models.SearchConditions{OrgId: orgId, AgentId: agentId}
	return crd.masterModel.DeleteByConditions(search)
}
