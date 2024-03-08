package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
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

	response.SuccessReturnData(ctx, "AuthController CheckUser name:"+name)
}

// SendSms 发送短信
func (u *AuthController) SendSms(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController SendSms")
}

// GetSliderCode 获取滑块码
func (u *AuthController) GetSliderCode(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AuthController GetSliderCode")
}

// Login 登录
func (u *AuthController) Login(ctx *gin.Context) {
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
