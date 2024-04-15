package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/utils/helpers"
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
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrganization + "name"))
	if helpers.StringIsEmpty(name) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	if admin_service.CreateAuthFactory().IsLocking(name) {
		response.FailSimplify(ctx, "用户已经被锁定，请稍后在试")
		return
	}

	adminModel := admin_service.CreateCurdFactory().GetItemByUsername(name)
	if adminModel.Id == 0 {
		response.FailSimplify(ctx, "用户名或密码错误")
		return
	}

	phoneCheck := 0
	if !helpers.StringIsEmpty(adminModel.Phone) && adminModel.LoginPhoneCheck == enums.IsYes {
		phoneCheck = 1
	}

	response.SuccessReturnData(ctx, gin.H{"phone_check": phoneCheck})
}

// SendSms 发送短信
func (u *AuthController) SendSms(ctx *gin.Context) {
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrganization + "name"))
	if helpers.StringIsEmpty(name) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	if admin_service.CreateAuthFactory().IsLocking(name) {
		response.FailSimplify(ctx, "用户已经被锁定，请稍后在试")
		return
	}

	adminModel := admin_service.CreateCurdFactory().GetItemByUsername(name)
	if adminModel.Id == 0 {
		response.FailSimplify(ctx, "用户名或密码错误")
		return
	}

	if helpers.StringIsEmpty(adminModel.Phone) {
		response.FailRoutine(ctx, -97, "没有添加手机号码")
		return
	}
	if adminModel.LoginPhoneCheck == enums.IsNot {
		response.FailRoutine(ctx, -98, "没有开启双因子验证")
		return
	}
	code := admin_service.CreateAuthFactory().GenerateBySms(adminModel.Phone)
	response.SuccessReturnData(ctx, gin.H{"code": code})
}

// GetSliderCode 获取滑块码
func (u *AuthController) GetSliderCode(ctx *gin.Context) {
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrganization + "name"))
	if helpers.StringIsEmpty(name) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	if admin_service.CreateAuthFactory().IsLocking(name) {
		response.FailSimplify(ctx, "用户已经被锁定，请稍后在试")
		return
	}
	code := admin_service.CreateAuthFactory().GenerateBySlider(name)
	response.SuccessReturnData(ctx, gin.H{"code": code})
}

// NormalLogin 登录
func (u *AuthController) NormalLogin(ctx *gin.Context) {
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrganization + "name"))
	pw := u.decrypt(ctx.GetString(consts.ValidatorPrefixOrganization + "pw"))
	sliderCode := ctx.GetString(consts.ValidatorPrefixOrganization + "sliderCode")
	smsCode := ctx.GetString(consts.ValidatorPrefixOrganization + "smsCode")
	if helpers.StringIsEmpty(name) || helpers.StringIsEmpty(pw) || helpers.StringIsEmpty(sliderCode) {
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

// HsLogin 幻势登录
func (u *AuthController) HsLogin(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController HsLogin")
}

// Logout 幻势登录
func (u *AuthController) Logout(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController Logout")
}
