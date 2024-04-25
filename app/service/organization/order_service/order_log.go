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

// Page 分页
func (crd *OrderLog) Page(orderNo string, offset, pageSize int, orderField, orderSort string) (total int64, items []*model_organization.OrderLogModel) {
	wheres := make(map[string]interface{})
	wheres["order_no"] = orderNo
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(wheres, offset, pageSize, order)
	return total, items
}
