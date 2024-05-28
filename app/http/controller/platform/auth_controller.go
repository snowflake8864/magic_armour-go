package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
)

// AuthController 安全验证相关
type AuthController struct {
}

func (u *AuthController) decrypt(data string) (result string) {
	result, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(data)
	return result
}

// CheckUser 验证用户
func (u *AuthController) CheckUser(ctx *gin.Context) {
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixPlatform + "name"))
	if string_helper.IsEmpty(name) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	authFact := admin_service.CreateAuthFactory()
	if err := authFact.IsLocking(name); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	err, adminModel := authFact.HasUser(name)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	phoneCheck := 0
	if !string_helper.IsEmpty(adminModel.Phone) && adminModel.LoginPhoneCheck == enums.IsYes {
		phoneCheck = 1
	}

	response.SuccessReturnData(ctx, gin.H{"phone_check": phoneCheck})
}

// SendSms 发送短信
func (u *AuthController) SendSms(ctx *gin.Context) {
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixPlatform + "name"))
	if string_helper.IsEmpty(name) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	authFact := admin_service.CreateAuthFactory()
	if err := authFact.IsLocking(name); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	err, adminModel := authFact.HasUser(name)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	if string_helper.IsEmpty(adminModel.Phone) {
		response.FailRoutine(ctx, -97, "没有添加手机号码")
		return
	}
	if adminModel.LoginPhoneCheck == enums.IsNot {
		response.FailRoutine(ctx, -98, "没有开启双因子验证")
		return
	}
	authFact.GenerateBySms(adminModel.Phone)

	response.SuccessSimplify(ctx)
}

// GetSliderCode 获取滑块码
func (u *AuthController) GetSliderCode(ctx *gin.Context) {
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixPlatform + "name"))
	if string_helper.IsEmpty(name) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	authFact := admin_service.CreateAuthFactory()
	if err := authFact.IsLocking(name); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	code := authFact.GenerateBySlider(name)
	response.SuccessReturnData(ctx, gin.H{"code": code})
}

// NormalLogin 登录
func (u *AuthController) NormalLogin(ctx *gin.Context) {
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixPlatform + "name"))
	pw := u.decrypt(ctx.GetString(consts.ValidatorPrefixPlatform + "pw"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixPlatform + "sliderCode")
	smsCode := ctx.GetString(consts.ValidatorPrefixPlatform + "smsCode")
	if string_helper.IsEmpty(name) || string_helper.IsEmpty(pw) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	result, errCode, err := admin_service.CreateAuthFactory().LoginByPw(name, pw, sliderCode, smsCode)
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessReturnData(ctx, result)
}

// Logout 退出
func (u *AuthController) Logout(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController Logout")
}
