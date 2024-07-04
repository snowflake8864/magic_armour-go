package order_cmd

import (
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/organization/order_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateCmdFactory() *OrderCmd {
	return &OrderCmd{}
}

type OrderCmd struct {
}

func (crd *OrderCmd) getOrgModels() []*organization_models.OrganizationModel {
	return organization_models.CreateOrganizationFactory().GetListByConditions(organization_models.SearchConditions{
		StatusList: []int{enum_organization.StatusNormal, enum_organization.StatusLocking},
	})
}

// UnpaidExpired  订单未支付且已过期的
func (crd *OrderCmd) UnpaidExpired() {
	variable.ZapLog.Sugar().Debug("UnpaidExpired start")
	yesterday, _, _, _ := date_helper.CalcByDay(-1)

	orgModels := crd.getOrgModels()
	if len(orgModels) > 0 {
		for _, model := range orgModels {
			order_service.CreateCurdFactory().AutoExpire(model.Id, yesterday)
		}
	}

	variable.ZapLog.Sugar().Debug("UnpaidExpired end yesterday:" + yesterday)
}
