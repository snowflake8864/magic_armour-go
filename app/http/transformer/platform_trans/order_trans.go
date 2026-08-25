package platform_trans

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/struct_helper"
)

type Order struct {
	transformer.Common
	Lang string
}

func (t *Order) pageGeneral(model *order_models.OrderModel) map[string]interface{} {
	data := struct_helper.SimpleStructToMapWithJsonTags(model)

	productGrantConfig := enum_product.ProductGrantConfig{}
	_ = json_helper.ToStruct(model.GrantModule, &productGrantConfig)
	data["grant_module"] = productGrantConfig
	data["pay_info"] = json_helper.ToMap(model.PayInfo)

	if t.Lang == consts.LanguageZh {
		data["grant_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type", model.GrantType, enum_product.EnumIntMap)
		data["grant_module_txt"] = productGrantConfig.SelectedTxt()
		data["grant_duration_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("grant_duration", model.GrantDuration, enum_order.EnumMap)
		data["order_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("order_status", model.OrderStatus, enum_order.EnumMap)
		data["is_expire_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("expire_status", model.IsExpire, enum_order.EnumMap)
		data["pay_method_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("pay_method", model.PayMethod, enum_order.EnumMap)
		data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source", model.Source, enums.EnumMapInt)
	} else {
		data["grant_type_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("type_en", model.GrantType, enum_product.EnumIntMap)
		data["grant_module_txt"] = productGrantConfig.SelectedTxtEn()
		data["grant_duration_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("grant_duration_en", model.GrantDuration, enum_order.EnumMap)
		data["order_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("order_status_en", model.OrderStatus, enum_order.EnumMap)
		data["is_expire_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("expire_status_en", model.IsExpire, enum_order.EnumMap)
		data["pay_method_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("pay_method_en", model.PayMethod, enum_order.EnumMap)
		data["source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source_en", model.Source, enums.EnumMapInt)
	}

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
			tmp := struct_helper.SimpleStructToMapWithJsonTags(model)
			if t.Lang == consts.LanguageZh {
				tmp["order_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("order_status", model.OrderStatus, enum_order.EnumMap)
				tmp["operate_action_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("action", model.OperateAction, enum_order.EnumMap)
				tmp["operate_source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source", model.OperateSource, enums.EnumMapInt)
			} else {
				tmp["order_status_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("order_status_en", model.OrderStatus, enum_order.EnumMap)
				tmp["operate_action_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("action_en", model.OperateAction, enum_order.EnumMap)
				tmp["operate_source_label"] = enums.Get1stIsStrAnd2ndIsIntReturnStr("operate_source_en", model.OperateSource, enums.EnumMapInt)
			}

			items = append(items, tmp)
		}
	}

	return gin.H{
		"total": total,
		"items": items,
	}
}
