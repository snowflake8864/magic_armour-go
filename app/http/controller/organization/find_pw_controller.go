package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/http/validator/common/common_regexp"
	"goskeleton/app/service/organization/user_service"
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
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if err := common_regexp.CheckPhone(userPhone); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	code := user_service.CreateAuthFactory().GenerateBySlider(user_service.ActonFindPw, userPhone)
	response.SuccessReturnData(ctx, gin.H{"code": code})
}

// SendSms 发送短信
func (u *FindPwController) SendSms(ctx *gin.Context) {
	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixOrg + "sliderCode")
	if string_helper.IsEmpty(userPhone) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	authFact := user_service.CreateAuthFactory()
	if !authFact.CheckBySlider(user_service.ActonFindPw, userPhone, sliderCode) {
		response.FailRoutine(ctx, -88, my_errors.CurdSliderIncorrect)
		return
	}

	if model := user_service.CreateCurdFactory().GetItemByPhone(userPhone); model == nil || model.Id == 0 {
		response.FailSimplify(ctx, my_errors.CurdUserNotExist)
		return
	}

	authFact.GenerateBySms(user_service.ActonFindPw, userPhone)

	response.SuccessSimplify(ctx)
}

// CheckSms 校验短信
func (u *FindPwController) CheckSms(ctx *gin.Context) {

	userPhone := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "userPhone"))
	smsCode := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "smsCode"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixOrg + "sliderCode")

	if string_helper.IsEmpty(userPhone) || string_helper.IsEmpty(smsCode) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if !user_service.CreateAuthFactory().CheckBySms(user_service.ActonFindPw, userPhone, smsCode) {
		response.FailRoutine(ctx, -99, my_errors.CurdSmsVerifyCodeNotMatch)
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
		response.FailSimplify(ctx, my_errors.CurdInvalidParameter)
		return
	}

	if !user_service.CreateAuthFactory().CheckBySms(user_service.ActonFindPw, userPhone, smsCode) {
		response.FailRoutine(ctx, -99, my_errors.CurdSmsVerifyCodeNotMatch)
		return
	} else {
		user_service.CreateAuthFactory().DeleteBySms(user_service.ActonFindPw, userPhone)
	}

	if err := user_service.CreateCurdFactory().FindPwByPhone(userPhone, userPw, ctx.ClientIP()); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}
