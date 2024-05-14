package order_service

import (
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/utils/helpers"
)

func CreateLogFactory() *OrderLog {
	return &OrderLog{masterModel: order_models.CreateOrderLogFactory()}
}

type OrderLog struct {
	OrderBase
	masterModel *order_models.OrderLogModel
}

// Add 新增订单日志
func (crd *OrderLog) Add(param *order_models.OrderLogModel) bool {
	if helpers.StringIsEmpty(param.CreatedAt) {
		param.CreatedAt = helpers.DatetimeCurr()
	}
	return crd.masterModel.InsertItem(param)
}

// DeleteByOrder 通过订单删除
func (crd *OrderLog) DeleteByOrder(orderNo string) bool {
	return crd.masterModel.DeleteByOrder(orderNo)
}

// Page 分页
func (crd *OrderLog) Page(orderNo string, offset, pageSize int, orderField, orderSort string) (total int64, items []*order_models.OrderLogModel) {
	wheres := make(map[string]interface{})
	wheres["order_no"] = orderNo
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(wheres, offset, pageSize, order)
	return total, items
}
