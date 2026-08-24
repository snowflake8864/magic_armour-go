package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/platform/prodcut_models"
	"goskeleton/app/service/platform/product_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/response"
)

// ProductController 终端相关
type ProductController struct {
	ControllerBase
}

// GrantTypes 授权类型
func (u *ProductController) GrantTypes(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("type", enum_product.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_product.EnumIntMap)
	if cur_userinfo.GetLocales(ctx) == consts.LanguageEn {
		labels = enums.Get1stIsStrReturnMapIntOfStr("type_en", enum_product.EnumIntMap)
	}
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatIntList(keys, labels, true))
}

// GrantModules 授权模块
func (u *ProductController) GrantModules(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfStr("module", enum_product.WebMapStr)
	labels := enums.Get1stIsStrReturnMapStrOfStr("module", enum_product.EnumStrMap)
	if cur_userinfo.GetLocales(ctx) == consts.LanguageEn {
		labels = enums.Get1stIsStrReturnMapStrOfStr("module_en", enum_product.EnumStrMap)
	}
	response.SuccessReturnData(ctx, (&transformer.Common{Lang: cur_userinfo.GetLocales(ctx)}).FormatStringList(keys, labels, false))
}

// AllPublished 全部已发布的产品
func (u *ProductController) AllPublished(ctx *gin.Context) {
	models := product_service.CreateCurdFactory().GetPublished()

	response.SuccessReturnData(ctx, (&platform_trans.Product{}).TransList(models))
}

// PageForManage 管理: 分页
func (u *ProductController) PageForManage(ctx *gin.Context) {
	//页面上传
	conditions := prodcut_models.SearchConditions{}
	conditions.Title = ctx.GetString(consts.ValidatorPrefixPlatform + "title")
	conditions.GrantType = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "grantType"))

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	//产品分页
	total, models := prodcut_models.CreateProductFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.Product{}).TransPage(total, models))
}

// AddForManage 管理: 添加
func (u *ProductController) AddForManage(ctx *gin.Context) {
	param := u.addParam(ctx)

	if !product_service.CreateCurdFactory().Add(&param) {
		response.FailSimplify(ctx, "新增失败")
		return
	}

	response.SuccessSimplify(ctx)
}

// EditForManage 管理: 编辑
func (u *ProductController) EditForManage(ctx *gin.Context) {
	param := u.addParam(ctx)

	if errCode, err := product_service.CreateCurdFactory().UpdateById(&param); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// DeleteForManage 管理: 删除
func (u *ProductController) DeleteForManage(ctx *gin.Context) {
	id := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "id"))
	if !product_service.CreateCurdFactory().DeleteById(id) {
		response.FailSimplify(ctx, my_errors.CurdDeleteFailed)
		return
	}

	response.SuccessSimplify(ctx)
}

func (u *ProductController) addParam(ctx *gin.Context) prodcut_models.ProductModel {
	var param = prodcut_models.ProductModel{}
	param.Id = int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "id"))
	param.GrantType = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "grant_type"))
	param.Title = ctx.GetString(consts.ValidatorPrefixPlatform + "title")
	param.Price = ctx.GetFloat64(consts.ValidatorPrefixPlatform + "price")
	param.Remark = ctx.GetString(consts.ValidatorPrefixPlatform + "remark")

	jsonStr := ctx.GetString(consts.ValidatorPrefixPlatform + "grant_module")
	productGrantConfig := enum_product.ProductGrantConfig{}
	if err := json_helper.ToStruct(jsonStr, &productGrantConfig); err != nil {
		productGrantConfig = enum_product.ProductGrantConfig{}
	}
	productGrantConfig.Default()

	param.GrantModule = json_helper.Encode(productGrantConfig)

	loginInfo := u.getLogin(ctx)
	param.CreateId = loginInfo.AdminId
	param.CreateName = loginInfo.Username
	return param
}
