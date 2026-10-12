package organization

import (
	"goskeleton/app/global/consts"
	"goskeleton/app/http/middleware/authorization/middleware_auth_org"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/service/platform/third_user_auth_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
	"net/url"

	"github.com/gin-gonic/gin"
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
	if err1, isPhone := authFact.IsPhoneCheck(userModel.Phone, userModel.OrgId); err1 != nil {
		response.FailSimplify(ctx, err1.Error())
		return
	} else if isPhone {
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
	authFact.GenerateBySms(user_service.ActonLogin, userModel.Phone)

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

	code := authFact.GenerateBySlider(user_service.ActonLogin, name)
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
	token := ctx.GetString(consts.ValidatorPrefixOrg + "token")

	deviceId, err := huanshi_platform.CreateCallbackFactory().CheckToken(token)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	backInfo, errCode, err := user_service.CreateAuthFactory().LoginByHs(deviceId, ctx.ClientIP())
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessReturnData(ctx, backInfo)
}

func (u *AuthController) getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	// 代理场景
	if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	return scheme + "://" + c.Request.Host + url.PathEscape("/#/login")
}

// KoalAddr 格尔登录地址
func (u *AuthController) KoalAddr(ctx *gin.Context) {
	//获取最新配置
	authConfig := third_user_auth_service.CreateConfigFactory().GetNewestItem()
	if authConfig == nil || string_helper.IsEmpty(authConfig.Addr) {
		response.FailSimplify(ctx, "auth config fail")
		return
	}

	response.SuccessReturnData(ctx, gin.H{
		"addr": third_user_auth_service.CreateRequestFactory(authConfig).GenerateLoginUrl(u.getBaseURL(ctx)),
	})

}

// KoalValidate 格尔验证票据
func (u *AuthController) KoalValidate(ctx *gin.Context) {
	ticket := ctx.GetString(consts.ValidatorPrefixOrg + "ticket")

	//获取最新配置
	authConfig := third_user_auth_service.CreateConfigFactory().GetNewestItem()
	if authConfig == nil || string_helper.IsEmpty(authConfig.Addr) {
		response.FailSimplify(ctx, "auth config fail")
		return
	}
	//
	err, ticketAttr := third_user_auth_service.CreateRequestFactory(authConfig).ValidateTicket(ticket, u.getBaseURL(ctx))
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	backInfo, errCode, err := user_service.CreateAuthFactory().LoginByKoal(string_helper.ToInt64(ticketAttr.Id), ctx.ClientIP())
	if err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessReturnData(ctx, backInfo)
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
