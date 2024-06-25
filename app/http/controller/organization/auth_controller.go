package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/http/middleware/authorization/middleware_auth_org"
	"goskeleton/app/service/organization/user_service"
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
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "name"))
	if string_helper.IsEmpty(name) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	authFact := user_service.CreateAuthFactory()
	if err := authFact.IsLocking(name); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	err, userModel := authFact.HasUser(name)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	phoneCheck := 0
	if err, smsSwitch := authFact.IsPhoneCheck(userModel.Phone, userModel.OrgId); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	} else if smsSwitch {
		phoneCheck = 1
	}

	response.SuccessReturnData(ctx, gin.H{"phone_check": phoneCheck})
}

// SendSms 发送短信
func (u *AuthController) SendSms(ctx *gin.Context) {
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "name"))
	if string_helper.IsEmpty(name) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	authFact := user_service.CreateAuthFactory()
	if err := authFact.IsLocking(name); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	err, userModel := authFact.HasUser(name)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	if string_helper.IsEmpty(userModel.Phone) {
		response.FailRoutine(ctx, -97, "没有开启双因子验证")
		return
	}
	if err, smsSwitch := authFact.IsPhoneCheck(userModel.Phone, userModel.OrgId); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	} else if !smsSwitch {
		response.FailRoutine(ctx, -98, "没有开启双因子验证")
		return
	}
	authFact.GenerateBySms(userModel.Phone)

	response.SuccessSimplify(ctx)
}

// GetSliderCode 获取滑块码
func (u *AuthController) GetSliderCode(ctx *gin.Context) {
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "name"))
	if string_helper.IsEmpty(name) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	authFact := user_service.CreateAuthFactory()
	if err := authFact.IsLocking(name); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	code := authFact.GenerateBySlider(name)
	response.SuccessReturnData(ctx, gin.H{"code": code})
}

// NormalLogin 登录
func (u *AuthController) NormalLogin(ctx *gin.Context) {
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "name"))
	pw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrg + "pw"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixOrg + "sliderCode")
	smsCode := ctx.GetString(consts.ValidatorPrefixOrg + "smsCode")
	if string_helper.IsEmpty(name) || string_helper.IsEmpty(pw) || string_helper.IsEmpty(sliderCode) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}

	result, errCode, err := user_service.CreateAuthFactory().LoginByPw(name, pw, sliderCode, smsCode, ctx.ClientIP())
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessReturnData(ctx, result)
}

// HsLogin 幻势登录
func (u *AuthController) HsLogin(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController HsLogin")
}

// Logout 幻势登录
func (u *AuthController) Logout(ctx *gin.Context) {
	headerParams := middleware_auth_org.Params{}
	if err := ctx.ShouldBindHeader(&headerParams); err != nil {
		response.SuccessSimplify(ctx)
		return
	}

	user_service.CreateAuthFactory().Logout(headerParams.Authorization)

	response.SuccessSimplify(ctx)
}
