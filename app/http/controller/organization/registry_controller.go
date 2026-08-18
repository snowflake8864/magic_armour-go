package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/http/transformer"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/model/organization/user_models"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
)

// RegistryController 注册相关
type RegistryController struct {
}

func (u *RegistryController) decrypt(data string) (result string) {
	result, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(data)
	return result
}

// Types 企业类型
func (u *RegistryController) Types(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("type", enum_organization.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("type", enum_organization.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// CheckOrgName 验证企业名称
func (u *RegistryController) CheckOrgName(ctx *gin.Context) {
	orgName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "orgName"))
	if string_helper.IsEmpty(orgName) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if model := organize_service.CreateOrgCurdFactory().GetItemByName(orgName); model != nil && model.Id > 0 {
		response.FailSimplify(ctx, my_errors.CurdOrgNameAlreadyExists)
		return
	}

	response.SuccessSimplify(ctx)
}

// CheckUserName 验证用户名称
func (u *RegistryController) CheckUserName(ctx *gin.Context) {
	userName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userName"))
	if string_helper.IsEmpty(userName) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if model := user_service.CreateCurdFactory().GetItemByName(userName); model != nil && model.Id > 0 {
		response.FailSimplify(ctx, my_errors.CurdUsernameAlreadyExists)
		return
	}

	response.SuccessSimplify(ctx)
}

// CheckUserPhone 验证用户手机
func (u *RegistryController) CheckUserPhone(ctx *gin.Context) {
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	if string_helper.IsEmpty(userPhone) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if model := user_service.CreateCurdFactory().GetItemByPhone(userPhone); model != nil && model.Id > 0 {
		response.FailSimplify(ctx, my_errors.CurdPhoneAlreadyExists)
		return
	}

	response.SuccessSimplify(ctx)
}

// GetSliderCode 获取滑块码
func (u *RegistryController) GetSliderCode(ctx *gin.Context) {
	userName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userName"))
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	if string_helper.IsEmpty(userName) || string_helper.IsEmpty(userPhone) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if model := user_service.CreateCurdFactory().GetItemByName(userName); model != nil && model.Id > 0 {
		response.FailSimplify(ctx, my_errors.CurdUsernameAlreadyExists)
		return
	}

	if model := user_service.CreateCurdFactory().GetItemByPhone(userPhone); model != nil && model.Id > 0 {
		response.FailSimplify(ctx, my_errors.CurdPhoneAlreadyExists)
		return
	}

	code := user_service.CreateAuthFactory().GenerateBySlider(user_service.ActonReg, userName)
	response.SuccessReturnData(ctx, gin.H{"code": code})
}

// SendSms 发送短信
func (u *RegistryController) SendSms(ctx *gin.Context) {
	userName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userName"))
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixOrg + "sliderCode")
	if string_helper.IsEmpty(userName) || string_helper.IsEmpty(userPhone) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	authFact := user_service.CreateAuthFactory()
	if !authFact.CheckBySlider(user_service.ActonReg, userName, sliderCode) {
		response.FailRoutine(ctx, -88, my_errors.CurdSliderIncorrect)
		return
	}

	authFact.GenerateBySms(user_service.ActonReg, userPhone)

	response.SuccessSimplify(ctx)
}

// Signup 注册
func (u *RegistryController) Signup(ctx *gin.Context) {
	orgName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "orgName"))
	userName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userName"))
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	userEmail := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userEmail"))
	userPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPw"))
	smsCode := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "smsCode"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixOrg + "sliderCode")
	orgType := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "orgType"))

	if string_helper.IsEmpty(orgName) || string_helper.IsEmpty(userName) || string_helper.IsEmpty(userPhone) ||
		string_helper.IsEmpty(userEmail) || string_helper.IsEmpty(userPw) ||
		string_helper.IsEmpty(smsCode) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if !user_service.CreateAuthFactory().CheckBySms(user_service.ActonReg, userPhone, smsCode) {
		response.FailRoutine(ctx, -99, my_errors.CurdSmsVerifyCodeNotMatch)
		return
	} else {
		user_service.CreateAuthFactory().DeleteBySms(user_service.ActonReg, userPhone)
	}

	//企业信息
	orgParam := organization_models.OrganizationModel{}
	orgParam.Name = orgName
	orgParam.UserName = userName
	orgParam.UserPhone = userPhone
	orgParam.UserEmail = userEmail
	orgParam.Type = orgType
	orgParam.Source = enums.OperateSourceOrganize

	//用户信息
	userParam := user_models.UserModel{}
	userParam.Username = orgParam.UserName
	userParam.Phone = orgParam.UserPhone
	userParam.Email = orgParam.UserEmail
	userParam.Password = userPw
	userParam.IsRoot = enums.IsYes

	errCode, err := organize_service.CreateOrgCurdFactory().Add(&orgParam, &userParam, ctx.ClientIP())
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
