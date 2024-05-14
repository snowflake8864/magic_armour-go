package access_log_service

import (
	"goskeleton/app/model/model_agent"
)

func CreateCurdFactory() *AccessLogCurd {
	return &AccessLogCurd{model_agent.CreateAccessLogFactory()}
}

type AccessLogCurd struct {
	masterModel *model_agent.AccessLogModel
}

// Add 新增
func (crd *AccessLogCurd) Add(param *model_agent.AccessLogModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
