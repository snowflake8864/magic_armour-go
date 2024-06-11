package platform_operate_service

import (
	"goskeleton/app/model/platform/operate_models"
)

func CreateCurdFactory() *OperateLogCurd {
	return &OperateLogCurd{masterModel: operate_models.CreateLogFactory()}
}

type OperateLogCurd struct {
	masterModel *operate_models.LogModel
}

// Add 新增
func (crd *OperateLogCurd) Add(param *operate_models.LogModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
