package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/service/organization/organize_service"
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

// PageForManage 管理: 分页
func (u *ProductController) PageForManage(ctx *gin.Context) {
	//查询条件
	conditions, err := u.getOrgSearchConditions(ctx)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	//默认状态
	conditions.StatusList = []int{
		enum_organization.StatusNormal,
		enum_organization.StatusLocking,
	}

	//分页参数和排序参数
	offset, pageSize, orderField, orderSort := u.getPageAndSort(ctx, enum_organization.OrgFieldDefault)

	//终端分页
	total, orgModels := organize_service.CreateOrgCurdFactory().Page(conditions, offset, pageSize, orderField, orderSort)

	response.SuccessReturnData(ctx, (&platform_trans.Org{}).TransPage(total, orgModels))
}
