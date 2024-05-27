package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
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

	response.SuccessSimplify(ctx)
}

// CheckOrgName 验证企业名称
func (u *RegistryController) CheckOrgName(ctx *gin.Context) {
	orgName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "orgName"))
	if string_helper.IsEmpty(orgName) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	response.SuccessSimplify(ctx)
}

// CheckUserName 验证用户名称
func (u *RegistryController) CheckUserName(ctx *gin.Context) {
	userName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userName"))
	if string_helper.IsEmpty(userName) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	response.SuccessSimplify(ctx)
}

// CheckUserPhone 验证用户手机
func (u *RegistryController) CheckUserPhone(ctx *gin.Context) {
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	if string_helper.IsEmpty(userPhone) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	response.SuccessSimplify(ctx)
}

// GetSliderCode 获取滑块码
func (u *RegistryController) GetSliderCode(ctx *gin.Context) {
	userName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userName"))
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	if string_helper.IsEmpty(userName) || string_helper.IsEmpty(userPhone) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	response.SuccessSimplify(ctx)
}

// SendSms 发送短信
func (u *RegistryController) SendSms(ctx *gin.Context) {
	userName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userName"))
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixOrg + "sliderCode")
	if string_helper.IsEmpty(userName) || string_helper.IsEmpty(userPhone) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
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

	if string_helper.IsEmpty(orgName) || string_helper.IsEmpty(userName) || string_helper.IsEmpty(userPhone) ||
		string_helper.IsEmpty(userEmail) || string_helper.IsEmpty(userPw) ||
		string_helper.IsEmpty(smsCode) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	response.SuccessSimplify(ctx)
}
