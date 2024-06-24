package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/platform_trans"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
)

// OrganizeController 终端相关
type OrganizeController struct {
	ControllerBase
}

func (u *OrganizeController) decrypt(data string) (result string) {
	result, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(data)
	return result
}

// Types 企业类型
func (u *OrganizeController) Types(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("type", enum_organization.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_organization.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// Detail 企业详情
func (u *OrganizeController) Detail(ctx *gin.Context) {
	orgId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))
	model := organize_service.CreateOrgCurdFactory().GetItemByOrg(orgId)
	if model == nil || model.Id == 0 {
		response.FailSimplify(ctx, "企业不存在")
		return
	}
	response.SuccessReturnData(ctx, (&platform_trans.Org{}).TransDetail(model))
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
	offset, pageSize, order := u.getPageAndSort(ctx, enum_organization.OrgFieldDefault)

	//终端分页
	total, orgModels := organization_models.CreateOrganizationFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.Org{}).TransPage(total, orgModels))
}

// LockForManage 管理：锁定
func (u *OrganizeController) LockForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orgId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))
	reason := ctx.GetString(consts.ValidatorPrefixPlatform + "reason")

	errCode, err := organize_service.CreateOrgCurdFactory().Lock(orgId, reason, ctx.ClientIP(), loginInfo)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// UnlockForManage 管理：解锁
func (u *OrganizeController) UnlockForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orgId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))
	reason := ctx.GetString(consts.ValidatorPrefixPlatform + "reason")

	errCode, err := organize_service.CreateOrgCurdFactory().Unlock(orgId, reason, ctx.ClientIP(), loginInfo)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// CleanForManage 管理：清理业务数据
func (u *OrganizeController) CleanForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orgId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))
	reason := ctx.GetString(consts.ValidatorPrefixPlatform + "reason")

	errCode, err := organize_service.CreateOrgCurdFactory().Clean(orgId, reason, ctx.ClientIP(), loginInfo)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetHsSwitchForManage 管理: 设置幻势开关
func (u *OrganizeController) SetHsSwitchForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orgId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))
	openSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "switch"))

	errCode, err := organize_service.CreateOrgCurdFactory().SetHsSwitch(orgId, openSwitch, ctx.ClientIP(), loginInfo)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// SetLoginSmsSwitchForManage 管理: 设置企业登录短信验证开关
func (u *OrganizeController) SetLoginSmsSwitchForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orgId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))
	openSwitch := int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "switch"))

	errCode, err := organize_service.CreateOrgCurdFactory().SetLoginSmsSwitch(orgId, openSwitch, ctx.ClientIP(), loginInfo)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
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
	offset, pageSize, order := u.getPageAndSort(ctx, enum_organization.OrgFieldDefault)

	//终端分页
	total, orgModels := organization_models.CreateOrganizationFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.Org{}).TransPage(total, orgModels))
}

// AddForVerify 审核: 添加
func (u *OrganizeController) AddForVerify(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	//企业信息
	orgParam := organization_models.OrganizationModel{}
	orgParam.Name = ctx.GetString(consts.ValidatorPrefixPlatform + "oraName")
	orgParam.UserName = ctx.GetString(consts.ValidatorPrefixPlatform + "userName")
	orgParam.UserPhone = ctx.GetString(consts.ValidatorPrefixPlatform + "userPhone")
	orgParam.UserEmail = ctx.GetString(consts.ValidatorPrefixPlatform + "userEmail")
	orgParam.Remark = ctx.GetString(consts.ValidatorPrefixPlatform + "remark")
	orgParam.Type = int(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgType"))
	orgParam.Source = enums.OperateSourcePlatform
	orgParam.CreateId = loginInfo.AdminId
	orgParam.CreateName = loginInfo.Username

	//用户信息
	userParam := user_models.UserModel{}
	userParam.Username = orgParam.UserName
	userParam.Phone = orgParam.UserPhone
	userParam.Email = orgParam.UserEmail
	userParam.Password = u.decrypt(ctx.GetString(consts.ValidatorPrefixPlatform + "userPw"))
	userParam.IsRoot = enums.IsYes
	if string_helper.IsEmpty(userParam.Password) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	errCode, err := organize_service.CreateOrgCurdFactory().Add(&orgParam, &userParam, ctx.ClientIP())
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// ApproveForVerify 审核: 批准
func (u *OrganizeController) ApproveForVerify(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orgId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))
	reason := ctx.GetString(consts.ValidatorPrefixPlatform + "reason")

	errCode, err := organize_service.CreateOrgCurdFactory().Approved(orgId, reason, ctx.ClientIP(), loginInfo)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// RefuseForVerify 审核: 不通过
func (u *OrganizeController) RefuseForVerify(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orgId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))
	reason := ctx.GetString(consts.ValidatorPrefixPlatform + "reason")

	errCode, err := organize_service.CreateOrgCurdFactory().Refused(orgId, reason, ctx.ClientIP(), loginInfo)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForUser 用户: 分页
func (u *OrganizeController) PageForUser(ctx *gin.Context) {
	orgId := int64(ctx.GetFloat64(consts.ValidatorPrefixPlatform + "orgId"))
	conditions := user_models.SearchConditions{OrgId: orgId}

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	//用户分页
	total, models := user_models.CreateUserFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&platform_trans.User{}).TransPage(total, models))
}

// PageForLog 日志: 分页
func (u *OrganizeController) PageForLog(ctx *gin.Context) {
	searchConditions := organization_models.SearchConditions{}
	searchConditions.LogOraName = ctx.GetString(consts.ValidatorPrefixPlatform + "oraName")
	searchConditions.UserName = ctx.GetString(consts.ValidatorPrefixPlatform + "userName")

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	//用户分页
	total, models := organization_models.CreateLogFactory().Page(searchConditions, offset, pageSize, order)
	response.SuccessReturnData(ctx, (&platform_trans.Org{}).TransLogPage(total, models))
}
