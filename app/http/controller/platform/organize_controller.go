package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/response"
)

// OrganizeController 终端相关
type OrganizeController struct {
	ControllerBase
}

// Types 企业类型
func (u *OrganizeController) Types(ctx *gin.Context) {
	list := enums.Get1stIsStrReturnMapIntOfStr("type", enum_organization.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatMapIntStr2List(list, true))
}

// PageForManage 管理: 分页
func (u *OrganizeController) PageForManage(ctx *gin.Context) {
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

	response.SuccessReturnData(ctx, (&platform_trans.Org{}).TransOrgManagePage(total, orgModels))
}

// PageForVerify 审核: 分页
func (u *OrganizeController) PageForVerify(ctx *gin.Context) {
	//查询条件
	conditions, err := u.getOrgSearchConditions(ctx)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	//默认状态
	conditions.StatusList = []int{
		enum_organization.StatusVerifying,
	}

	//分页参数和排序参数
	offset, pageSize, orderField, orderSort := u.getPageAndSort(ctx, enum_organization.OrgFieldDefault)

	//终端分页
	total, orgModels := organize_service.CreateOrgCurdFactory().Page(conditions, offset, pageSize, orderField, orderSort)

	response.SuccessReturnData(ctx, (&platform_trans.Org{}).TransOrgManagePage(total, orgModels))
}
