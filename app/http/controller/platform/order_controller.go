package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/service/organization/order_service"
	"goskeleton/app/utils/response"
)

// OrderController 终端相关
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

// AddForOrg  为企业添加订单授权
func (u *OrderController) AddForOrg(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := order_models.OrderModel{}
	param.OrgId = int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))
	param.ProductId = int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "productId"))
	param.BuyNumber = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "buyNumber"))
	param.GrantDuration = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "grantDuration"))
	param.RealAmount = ctx.GetFloat64(consts.ValidatorPrefixPlatform + "realAmount")
	param.DeclineRemark = ctx.GetString(consts.ValidatorPrefixPlatform + "remark")
	param.Source = enums.OperateSourcePlatform

	if errCode, err := order_service.CreateCurdFactory().Add(&param, loginInfo.AdminId, loginInfo.Username); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// ValidByOrg 企业下有效订单列表
func (u *OrderController) ValidByOrg(ctx *gin.Context) {
	OrgId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))

	models := order_service.CreateCurdFactory().GetValidItemsByOrg(OrgId)

	response.SuccessReturnData(ctx, (&platform_trans.Order{}).TransList(models))
}

// PageForGrant 授权管理: 分页
func (u *OrderController) PageForGrant(ctx *gin.Context) {
	//查询条件
	conditions := u.getOrderSearchConditions(ctx)
	conditions.IsExpire = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "isExpire"))
	conditions.OrderStatus = enum_order.OrderStatusCompleted

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, enum_organization.OrgFieldDefault)

	//终端分页
	total, orderModels := order_models.CreateOrderFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.Order{}).TransPage(total, orderModels))
}

func (u *OrderController) getOrderSearchConditions(ctx *gin.Context) order_models.SearchConditions {
	searchConditions := order_models.SearchConditions{}
	searchConditions.OraName = ctx.GetString(consts.ValidatorPrefixPlatform + "oraName")
	searchConditions.UserName = ctx.GetString(consts.ValidatorPrefixPlatform + "userName")
	searchConditions.ProductName = ctx.GetString(consts.ValidatorPrefixPlatform + "productName")
	searchConditions.GteCreateAt = ctx.GetString(consts.ValidatorPrefixPlatform + "startTime")
	searchConditions.LteCreateAt = ctx.GetString(consts.ValidatorPrefixPlatform + "endTime")
	return searchConditions
}

// PageForManage 订单管理: 分页
func (u *OrderController) PageForManage(ctx *gin.Context) {
	//查询条件
	conditions := u.getOrderSearchConditions(ctx)
	conditions.IsExpire = enums.AllDefaultInt
	conditions.OrderStatus = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orderStatus"))
	conditions.GteCreateAt = ctx.GetString(consts.ValidatorPrefixPlatform + "startTime")
	conditions.LteCreateAt = ctx.GetString(consts.ValidatorPrefixPlatform + "endTime")

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, enum_organization.OrgFieldDefault)

	//终端分页
	total, orderModels := order_models.CreateOrderFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.Order{}).TransPage(total, orderModels))
}

// CancelForManage 订单管理: 取消
func (u *OrderController) CancelForManage(ctx *gin.Context) {
	orderNo := ctx.GetString(consts.ValidatorPrefixPlatform + "orderNo")

	loginInfo := u.getLogin(ctx)
	errCode, err := order_service.CreateCurdFactory().PlatformCancel(orderNo, loginInfo.AdminId, loginInfo.Username)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeclineForManage 订单管理: 直降
func (u *OrderController) DeclineForManage(ctx *gin.Context) {
	orderNo := ctx.GetString(consts.ValidatorPrefixPlatform + "orderNo")
	remark := ctx.GetString(consts.ValidatorPrefixPlatform + "remark")
	realAmount := ctx.GetFloat64(consts.ValidatorPrefixPlatform + "realAmount")

	loginInfo := u.getLogin(ctx)
	errCode, err := order_service.CreateCurdFactory().Decline(orderNo, remark, realAmount, loginInfo.AdminId, loginInfo.Username)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PaidForManage 订单管理: 确认已支付
func (u *OrderController) PaidForManage(ctx *gin.Context) {
	orderNo := ctx.GetString(consts.ValidatorPrefixPlatform + "orderNo")

	loginInfo := u.getLogin(ctx)
	errCode, err := order_service.CreateCurdFactory().PlatformCancel(orderNo, loginInfo.AdminId, loginInfo.Username)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// UnpaidForManage 订单管理: 确认未付款
func (u *OrderController) UnpaidForManage(ctx *gin.Context) {
	orderNo := ctx.GetString(consts.ValidatorPrefixPlatform + "orderNo")

	loginInfo := u.getLogin(ctx)
	errCode, err := order_service.CreateCurdFactory().PlatformCancel(orderNo, loginInfo.AdminId, loginInfo.Username)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DetailForManage 订单管理: 详情
func (u *OrderController) DetailForManage(ctx *gin.Context) {
	orderNo := ctx.GetString(consts.ValidatorPrefixPlatform + "orderNo")

	model := order_service.CreateCurdFactory().GetItemByOrderNo(orderNo)
	response.SuccessReturnData(ctx, (&platform_trans.Order{}).TransDetail(model))
}

// LogsForManage 订单管理:  日志列表
func (u *OrderController) LogsForManage(ctx *gin.Context) {
	orderNo := ctx.GetString(consts.ValidatorPrefixPlatform + "orderNo")
	searchConditions := order_models.SearchConditions{}
	searchConditions.OrderNo = orderNo
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, enum_organization.OrgFieldDefault)

	//日志分页
	total, logModels := order_models.CreateOrderLogFactory().Page(searchConditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.Order{}).TransLogsPage(total, logModels))
}
