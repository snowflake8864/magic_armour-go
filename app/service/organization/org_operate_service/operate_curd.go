package org_operate_service

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/model/organization/org_operate_models"
)

func CreateCurdFactory() *OperateLogCurd {
	return &OperateLogCurd{masterModel: org_operate_models.CreateLogFactory()}
}

type OperateLogCurd struct {
	masterModel *org_operate_models.LogModel
}

// Add 新增
func (crd *OperateLogCurd) Add(param *org_operate_models.LogModel) bool {
	param.Id = variable.SnowFlake.GetId()
	isOk := crd.masterModel.InsertItem(param)

	return isOk
}
