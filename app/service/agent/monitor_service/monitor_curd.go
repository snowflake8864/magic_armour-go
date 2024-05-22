package monitor_service

import (
	"goskeleton/app/model/agent/agent_monitor_models"
	"goskeleton/app/model/agent/port_models"
)

func CreateMonitorFactory() *MonitorCurd {
	return &MonitorCurd{agent_monitor_models.CreateMonitorFactory()}
}

type MonitorCurd struct {
	masterModel *agent_monitor_models.MonitorModel
}

// Add 新增
func (crd *MonitorCurd) Add(orgId int64, agentId string, param *agent_monitor_models.MonitorModel) bool {

	return crd.masterModel.InsertBatch(params)
}

// DelByOrg 删除：企业
func (crd *MonitorCurd) DelByOrg(orgId int64) bool {
	search := port_models.SearchConditions{OrgId: orgId}
	return crd.masterModel.DeleteByConditions(search)
}

// DelByOrgAndAgent 删除：企业和终端
func (crd *MonitorCurd) DelByOrgAndAgent(orgId int64, agentId string) bool {
	search := port_models.SearchConditions{OrgId: orgId, AgentId: agentId}
	return crd.masterModel.DeleteByConditions(search)
}
