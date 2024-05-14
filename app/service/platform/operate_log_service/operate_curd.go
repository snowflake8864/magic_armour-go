package operate_log_service

import (
	"goskeleton/app/model/platform/operate_models"
)

func CreateCurdFactory() *OperateLogCurd {
	return &OperateLogCurd{masterModel: operate_models.CreateOperateLogFactory()}
}

type OperateLogCurd struct {
	masterModel *operate_models.OperateLogModel
}

// Add 新增
func (crd *OperateLogCurd) Add(param *operate_models.OperateLogModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
