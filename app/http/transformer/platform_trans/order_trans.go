package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/utils/helpers/json_helper"
)

type Order struct {
	transformer.Common
}

func (t *Order) pageGeneral(model *order_models.OrderModel) map[string]interface{} {
	data := make(map[string]interface{})
	data["order_no"] = model.OrderNo
	data["org_id"] = model.OrgId
	data["org_name"] = model.OrgName
	data["user_email"] = model.UserEmail
	data["user_phone"] = model.UserPhone
	data["product_id"] = model.ProductId
	data["product_title"] = model.ProductTitle
	data["grant_type"] = model.GrantType
	data["grant_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.GrantType, enum_product.EnumIntMap)
	productGrantConfig := enum_product.ProductGrantConfig{}
	_ = json_helper.Json2Struct(model.GrantModule, &productGrantConfig)
	data["grant_module"] = productGrantConfig
	data["grant_module_txt"] = productGrantConfig.SelectedTxt()
	data["sale_price"] = model.SalePrice
	data["grant_duration"] = model.GrantDuration
	data["buy_number"] = model.BuyNumber
	data["usageNumber"] = model.UsageNumber
	data["total_amount"] = model.TotalAmount
	data["decline_amount"] = model.DeclineAmount
	data["real_amount"] = model.RealAmount
	data["order_status"] = model.OrderStatus
	data["order_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("order_status", model.OrderStatus, enum_order.EnumMap)
	data["is_expire"] = model.IsExpire
	data["is_expire_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("expire_status", model.IsExpire, enum_order.EnumMap)
	data["expire_at"] = model.ExpireAt
	data["pay_method"] = model.PayMethod
	data["pay_method_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("pay_method", model.PayMethod, enum_order.EnumMap)
	data["pay_info"] = json_helper.ToMap(model.PayInfo)
	data["pay_order_no"] = model.PayOrderNo
	data["created_at"] = model.CreatedAt
	data["decline_at"] = model.DeclineAt
	data["decline_remark"] = model.DeclineRemark
	data["pay_at"] = model.PayAt
	data["complete_at"] = model.CompleteAt
	data["cancel_at"] = model.CancelAt
	data["source"] = model.Source
	data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source", model.Source, enums.EnumMapInt)
	return data
}

// TransPage 转换成管理页需要的数据
func (t *Order) TransPage(total int64, models []*order_models.OrderModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}

// TransList 转换列表页需要的数据
func (t *Order) TransList(models []*order_models.OrderModel) (items []map[string]interface{}) {
	if models != nil {
		for _, model := range models {
			tmp := t.pageGeneral(model)
			items = append(items, tmp)
		}
	}
	return items
}

// TransDetail 转换详情页需要的数据
func (t *Order) TransDetail(model *order_models.OrderModel) (item map[string]interface{}) {
	if model != nil {
		item = t.pageGeneral(model)
	}
	return item
}

// TransLogsPage 转换成日志页需要的数据
func (t *Order) TransLogsPage(total int64, models []*order_models.OrderLogModel) interface{} {
	var items []map[string]interface{}
	if models != nil {
		for _, model := range models {
			tmp := make(map[string]interface{})

			tmp["org_id"] = model.OrgId
			tmp["OrderStatus"] = model.OrderStatus
			tmp["order_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("order_status", model.OrderStatus, enum_order.EnumMap)
			tmp["operate_id"] = model.OperateId
			tmp["operate_name"] = model.OperateName
			tmp["operate_action"] = model.OperateAction
			tmp["operate_action_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("action", model.OperateAction, enum_order.EnumMap)
			tmp["operate_source"] = model.OperateSource
			tmp["operate_source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source", model.OperateSource, enums.EnumMapInt)
			tmp["usage_number"] = model.UsageNumber
			tmp["usage_agents"] = model.UsageAgents
			tmp["remark"] = model.Remark
			tmp["created_at"] = model.CreatedAt
			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
