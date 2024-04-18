package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/model_organization"
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

// AddForVerify 审核: 分页
func (u *OrganizeController) AddForVerify(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	//企业信息
	orgParam := model_organization.OrganizationModel{}
	orgParam.Name = ctx.GetString(consts.ValidatorPrefixPlatform + "oraName")
	orgParam.UserName = ctx.GetString(consts.ValidatorPrefixPlatform + "userName")
	orgParam.UserPhone = ctx.GetString(consts.ValidatorPrefixPlatform + "userPhone")
	orgParam.UserEmail = ctx.GetString(consts.ValidatorPrefixPlatform + "userEmail")
	orgParam.Remark = ctx.GetString(consts.ValidatorPrefixPlatform + "remark")
	orgParam.Type = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgType"))
	orgParam.Source = enum_organization.SourcePlatform
	orgParam.CreateId = loginInfo.AdminId
	orgParam.CreateName = loginInfo.Username

	//用户信息
	userParam := model_organization.UserModel{}
	userParam.Username = orgParam.UserName
	userParam.Phone = orgParam.UserPhone
	userParam.Email = orgParam.UserEmail
	userParam.Password = ctx.GetString(consts.ValidatorPrefixPlatform + "userPw")
	userParam.IsRoot = enums.IsYes

	errCode, err := organize_service.CreateOrgCurdFactory().Add(&orgParam, &userParam, ctx.ClientIP())
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
