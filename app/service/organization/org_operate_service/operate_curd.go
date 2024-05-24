package org_operate_service

import (
	"goskeleton/app/model/organization/org_operate_models"
)

func CreateCurdFactory() *OperateLogCurd {
	return &OperateLogCurd{masterModel: org_operate_models.CreateOperateLogFactory()}
}

type OperateLogCurd struct {
	masterModel *org_operate_models.OperateLogModel
}

// Add 新增
func (crd *OperateLogCurd) Add(param *org_operate_models.OperateLogModel) bool {
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
