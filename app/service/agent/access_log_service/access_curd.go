package access_log_service

import (
	"goskeleton/app/model/agent/agent_access_models"
)

func CreateCurdFactory() *AccessLogCurd {
	return &AccessLogCurd{agent_access_models.CreateAccessLogFactory()}
}

type AccessLogCurd struct {
	masterModel *agent_access_models.AccessLogModel
}

// Add 新增
func (crd *AccessLogCurd) Add(param *agent_access_models.AccessLogModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
