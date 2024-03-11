package platform

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
	name := u.decrypt(ctx.GetString(consts.ValidatorPrefix + "name"))
	if helpers.StringIsEmpty(name) {
		response.FailSimplify(ctx, "参数不正确")
		return
	}
	adminModel := admin_service.CreateCurdFactory().GetItemByUsername(name)
	if adminModel.Id == 0 {
		response.FailSimplify(ctx, "用户名或密码错误")
		return
	}
	phoneCheck := 0
	if !helpers.StringIsEmpty(adminModel.Phone) && adminModel.PhoneLoginCheck == enums.IsYes {
		phoneCheck = 1
	}

	response.SuccessReturnData(ctx, gin.H{"phone_check": phoneCheck})
}

// SendSms 发送短信
func (u *AuthController) SendSms(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController SendSms")
}

// GetSliderCode 获取滑块码
func (u *AuthController) GetSliderCode(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController GetSliderCode")
}

// NormalLogin 登录
func (u *AuthController) NormalLogin(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController Login")
}

// HsLogin 幻势登录
func (u *AuthController) HsLogin(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController HsLogin")
}

// Logout 幻势登录
func (u *AuthController) Logout(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController Logout")
}
