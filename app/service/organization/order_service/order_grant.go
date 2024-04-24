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

func (crd *OrderGrant) Init(orderNo string, orgId int64, buyNumber int) (param model_organization.OrderGrantModel) {
	param.OrderNo = orderNo
	param.OrgId = orgId
	param.AgentTotal = buyNumber
	param.AgentUsed = 0
	return param
}
