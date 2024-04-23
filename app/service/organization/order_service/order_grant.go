package order_service

import (
	"goskeleton/app/model/model_organization"
)

func CreateGrantFactory() *OrderGrant {
	return &OrderGrant{masterModel: model_organization.CreateOrderGrantFactory()}
}

type OrderGrant struct {
	OrderBase
	masterModel *model_organization.OrderGrantModel
}
