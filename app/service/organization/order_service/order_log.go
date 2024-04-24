package order_service

import (
	"goskeleton/app/model/model_organization"
	"goskeleton/app/utils/helpers"
)

func CreateLogFactory() *OrderLog {
	return &OrderLog{masterModel: model_organization.CreateOrderLogFactory()}
}

type OrderLog struct {
	OrderBase
	masterModel *model_organization.OrderLogModel
}

// Add 新增订单日志
func (crd *OrderLog) Add(param *model_organization.OrderLogModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	return crd.masterModel.InsertItem(param)
}
