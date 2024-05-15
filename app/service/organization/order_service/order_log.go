package order_service

import (
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
)

func CreateLogFactory() *OrderLog {
	return &OrderLog{masterModel: order_models.CreateOrderLogFactory()}
}

type OrderLog struct {
	masterModel *order_models.OrderLogModel
}

// Add 新增订单日志
func (crd *OrderLog) Add(param *order_models.OrderLogModel) bool {
	if string_helper.IsEmpty(param.CreatedAt) {
		param.CreatedAt = date_helper.CurDatetime()
	}
	return crd.masterModel.InsertItem(param)
}

// DeleteByOrder 通过订单删除
func (crd *OrderLog) DeleteByOrder(orderNo string) bool {
	return crd.masterModel.DeleteByOrder(orderNo)
}
