package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/http/transformer"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/utils/response"
)

// OrderController  订单授权相关
type OrderController struct {
	ControllerBase
}

// OrderStatus  订单状态
func (u *OrderController) OrderStatus(ctx *gin.Context) {

	keys := enums.Get1stIsStrReturnSliceOfInt("order_status", enum_order.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("order_status", enum_order.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// ExpireStatus 过期状态
func (u *OrderController) ExpireStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("expire_status", enum_order.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("expire_status", enum_order.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// GrantDuration 授权时长
func (u *OrderController) GrantDuration(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("grant_duration", enum_order.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("grant_duration", enum_order.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// PayMethod 支付方式
func (u *OrderController) PayMethod(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("pay_method", enum_order.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("pay_method", enum_order.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// PayBankcard 支付银行卡信息
func (u *OrderController) PayBankcard(ctx *gin.Context) {

	maps := dictionary_service.CreateCurdFactory().GetMapByModule(enum_dictionary.ModuleOrder)

	result := map[string]string{
		"card_name": maps[enum_dictionary.OrderBankcardName],
		"card_no":   maps[enum_dictionary.OrderBankcardNo],
		"bank_name": maps[enum_dictionary.OrderBankcardBank],
		"tel":       maps[enum_dictionary.OrderPayingTel],
	}

	response.SuccessReturnData(ctx, result)
}

// ValidList 有效订单列表
func (u *OrderController) ValidList(ctx *gin.Context) {

	maps := dictionary_service.CreateCurdFactory().GetMapByModule(enum_dictionary.ModuleOrder)

	result := map[string]string{
		"card_name": maps[enum_dictionary.OrderBankcardName],
		"card_no":   maps[enum_dictionary.OrderBankcardNo],
		"bank_name": maps[enum_dictionary.OrderBankcardBank],
		"tel":       maps[enum_dictionary.OrderPayingTel],
	}

	response.SuccessReturnData(ctx, result)
}

// PageForPurchase 购买：分页
func (u *OrderController) PageForPurchase(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}

// InfoForPurchase 购买：信息
func (u *OrderController) InfoForPurchase(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}

// LogsForPurchase 购买：日志
func (u *OrderController) LogsForPurchase(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}

// AddForPurchase 购买：添加
func (u *OrderController) AddForPurchase(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}

// PayByBankcardForPurchase 购买：银卡支付
func (u *OrderController) PayByBankcardForPurchase(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}

// CancelForPurchase 购买：取消
func (u *OrderController) CancelForPurchase(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}

// DeleteForPurchase 购买：删除
func (u *OrderController) DeleteForPurchase(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}

// PageForGrant 终端授权：分页
func (u *OrderController) PageForGrant(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}

// SaveForGrant 终端授权：保存
func (u *OrderController) SaveForGrant(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}

// PageForUsed 已使用：分页
func (u *OrderController) PageForUsed(ctx *gin.Context) {

	response.SuccessSimplify(ctx)
}
