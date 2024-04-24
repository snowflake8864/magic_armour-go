package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/model_organization"
	"goskeleton/app/service/organization/order_service"
	"goskeleton/app/service/platform/product_service"
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

// PageForManage 管理: 分页
func (u *OrderController) PageForManage(ctx *gin.Context) {
	//页面上传
	conditions := product_service.SearchConditions{}
	conditions.Title = ctx.GetString(consts.ValidatorPrefixPlatform + "title")
	conditions.GrantType = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "grantType"))
	conditions.IsPublish = enums.AllDefaultInt

	//分页参数和排序参数
	offset, pageSize, orderField, orderSort := u.getPageAndSort(ctx, enum_organization.OrgFieldDefault)

	//产品分页
	total, models := product_service.CreateCurdFactory().Page(conditions, offset, pageSize, orderField, orderSort)

	response.SuccessReturnData(ctx, (&platform_trans.Product{}).TransPage(total, models))
}

// AddForOrg  为企业添加订单授权
func (u *OrderController) AddForOrg(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := model_organization.OrderModel{}
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
	offset, pageSize, orderField, orderSort := u.getPageAndSort(ctx, enum_organization.OrgFieldDefault)

	//终端分页
	total, orderModels := order_service.CreateCurdFactory().Page(conditions, offset, pageSize, orderField, orderSort)

	response.SuccessReturnData(ctx, (&platform_trans.Order{}).TransPage(total, orderModels))
}

func (u *OrderController) getOrderSearchConditions(ctx *gin.Context) order_service.SearchConditions {
	searchConditions := order_service.SearchConditions{}
	searchConditions.OraName = ctx.GetString(consts.ValidatorPrefixPlatform + "oraName")
	searchConditions.UserName = ctx.GetString(consts.ValidatorPrefixPlatform + "userName")
	searchConditions.ProductName = ctx.GetString(consts.ValidatorPrefixPlatform + "productName")
	searchConditions.CreateAtGte = ctx.GetString(consts.ValidatorPrefixPlatform + "startTime")
	searchConditions.CreateAtLte = ctx.GetString(consts.ValidatorPrefixPlatform + "endTime")
	return searchConditions
}
