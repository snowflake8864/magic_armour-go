package operate_log_service

import (
	"goskeleton/app/model/model_platform"
)

func CreateCurdFactory() *OperateLogCurd {
	return &OperateLogCurd{masterModel: model_platform.CreateOperateLogFactory()}
}

type OperateLogCurd struct {
	masterModel *model_platform.OperateLogModel
}

// Add 新增
func (crd *OperateLogCurd) Add(param *model_platform.OperateLogModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
