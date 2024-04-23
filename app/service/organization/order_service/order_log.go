package order_service

import (
	"goskeleton/app/model/model_organization"
)

func CreateLogFactory() *OrderLog {
	return &OrderLog{masterModel: model_organization.CreateOrderLogFactory()}
}

type OrderLog struct {
	OrderBase
	masterModel *model_organization.OrderLogModel
}

