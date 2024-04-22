package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/service/platform/product_service"
	"goskeleton/app/utils/response"
)

// ProductController 终端相关
type ProductController struct {
	ControllerBase
}

// GrantTypes 授权类型
func (u *ProductController) GrantTypes(ctx *gin.Context) {
	list := enums.Get1stIsStrReturnMapIntOfStr("type", enum_product.EnumIntMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatMapIntOfStr2List(list, true))
}

// GrantModules 授权模块
func (u *ProductController) GrantModules(ctx *gin.Context) {
	list := enums.Get1stIsStrReturnMapStrOfStr("module", enum_product.EnumStrMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatMapStrOStr2List(list, false))
}

// AllPublished 全部已发布的产品
func (u *ProductController) AllPublished(ctx *gin.Context) {
	models := product_service.CreateCurdFactory().GetPublished()

	response.SuccessReturnData(ctx, (&platform_trans.Product{}).TransList(models))
}

// PageForManage 管理: 分页
func (u *ProductController) PageForManage(ctx *gin.Context) {
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
