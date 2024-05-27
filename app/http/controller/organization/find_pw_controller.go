package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
)

// FindPwController 注册相关
type FindPwController struct {
}

func (u *FindPwController) decrypt(data string) (result string) {
	result, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(data)
	return result
}

// GetSliderCode 获取滑块码
func (u *FindPwController) GetSliderCode(ctx *gin.Context) {
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	if string_helper.IsEmpty(userPhone) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	response.SuccessSimplify(ctx)
}

// CheckOrgName 验证企业名称
func (u *FindPwController) CheckOrgName(ctx *gin.Context) {
	orgName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "orgName"))
	if string_helper.IsEmpty(orgName) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	response.SuccessSimplify(ctx)
}

// CheckUserName 验证用户名称
func (u *FindPwController) CheckUserName(ctx *gin.Context) {
	userName := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userName"))
	if string_helper.IsEmpty(userName) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	response.SuccessSimplify(ctx)
}

// CheckUserPhone 验证用户手机
func (u *FindPwController) CheckUserPhone(ctx *gin.Context) {
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	if string_helper.IsEmpty(userPhone) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	response.SuccessSimplify(ctx)
}

// SendSms 发送短信
func (u *FindPwController) SendSms(ctx *gin.Context) {
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixOrg + "sliderCode")
	if string_helper.IsEmpty(userPhone) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	response.SuccessSimplify(ctx)
}

// CheckSms 校验短信
func (u *FindPwController) CheckSms(ctx *gin.Context) {

	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	smsCode := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "smsCode"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixOrg + "sliderCode")

	if string_helper.IsEmpty(userPhone) || string_helper.IsEmpty(smsCode) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	response.SuccessSimplify(ctx)
}

// ModifyPw 修改密码
func (u *FindPwController) ModifyPw(ctx *gin.Context) {

	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	userPw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPw"))
	smsCode := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "smsCode"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixOrg + "sliderCode")

	if string_helper.IsEmpty(userPhone) || string_helper.IsEmpty(userPw) || string_helper.IsEmpty(smsCode) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	response.SuccessSimplify(ctx)
}
