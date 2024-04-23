package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/service/platform/product_service"
	"goskeleton/app/utils/helpers"
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

// AddForManage 管理: 添加
func (u *OrderController) AddForManage(ctx *gin.Context) {
	param := u.addParam(ctx)

	if !product_service.CreateCurdFactory().Add(&param) {
		response.FailSimplify(ctx, "新增失败")
		return
	}

	response.SuccessSimplify(ctx)
}

// EditForManage 管理: 编辑
func (u *OrderController) EditForManage(ctx *gin.Context) {
	param := u.addParam(ctx)

	if errCode, err := product_service.CreateCurdFactory().UpdateById(&param); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// DeleteForManage 管理: 删除
func (u *OrderController) DeleteForManage(ctx *gin.Context) {
	id := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "id"))
	if !product_service.CreateCurdFactory().DeleteById(id) {
		response.FailSimplify(ctx, "删除失败")
		return
	}

	response.SuccessSimplify(ctx)
}

func (u *OrderController) addParam(ctx *gin.Context) model_platform.ProductModel {
	var param = model_platform.ProductModel{}
	param.Id = int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "id"))
	param.GrantType = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "grant_type"))
	param.Title = ctx.GetString(consts.ValidatorPrefixPlatform + "title")
	param.Price = ctx.GetFloat64(consts.ValidatorPrefixPlatform + "price")
	param.Remark = ctx.GetString(consts.ValidatorPrefixPlatform + "remark")

	productGrantConfig := enum_product.ProductGrantConfig{}
	productGrantConfig.Process = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "m_process"))
	productGrantConfig.Tamper = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "m_tamper"))
	productGrantConfig.Blackmail = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "m_blackmail"))
	productGrantConfig.Port = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "m_port"))
	productGrantConfig.Plug = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "m_plug"))
	productGrantConfig.Usb = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "m_usb"))
	productGrantConfig.Syslog = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "m_syslog"))
	productGrantConfig.Outreach = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "m_outreach"))
	productGrantConfig.Login = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "m_login"))
	param.GrantModule = helpers.JsonEncode(productGrantConfig)

	loginInfo := u.getLogin(ctx)
	param.CreateId = loginInfo.AdminId
	param.CreateName = loginInfo.Username
	return param
}
