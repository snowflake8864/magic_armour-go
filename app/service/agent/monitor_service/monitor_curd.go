package monitor_service

import (
	"errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_monitor_models"
)

func CreateMonitorFactory() *MonitorCurd {
	return &MonitorCurd{agent_monitor_models.CreateMonitorFactory()}
}

type MonitorCurd struct {
	masterModel *agent_monitor_models.MonitorModel
}

// Add 新增
func (crd *MonitorCurd) Add(param *agent_monitor_models.MonitorModel) error {
	if param.Id == 0 {
		param.Id = variable.SnowFlake.GetId()
	}
	if crd.masterModel.InsertItem(param) {
		return nil
	}

	return errors.New("保存失败")
}

// DelByOrg 删除：企业
func (crd *MonitorCurd) DelByOrg(orgId int64) bool {
	return crd.masterModel.DeleteByOrg(orgId)
}

// DelByOrgAndAgent 删除：企业和终端
func (crd *MonitorCurd) DelByOrgAndAgent(orgId int64, agentId string) bool {
	return crd.masterModel.DeleteByOrgAndAgent(orgId, agentId)
}
