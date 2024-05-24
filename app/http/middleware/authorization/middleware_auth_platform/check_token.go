package middleware_auth_platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/utils/response"
)

type Params struct {
	Authorization string `header:"Authorization" binding:"required,min=32,max=64"`
}

// CheckToken 检查平台token完整性、有效性中间件
func CheckToken() gin.HandlerFunc {
	//variable.ZapLog.Info("检查终端token完整性、有效性中间件----------------")
	return func(context *gin.Context) {
		headerParams := Params{}
		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.ErrorTokenFail(context, consts.PlatformLoginExpiredCode, consts.JwtTokenMustValid+err.Error())
			return
		}
		//variable.ZapLog.Info("Authorization:" + headerParams.Authorization)
		loginInfo, err := admin_service.CreateAuthFactory().DecryptByLogin(headerParams.Authorization, true)
		// token验证通过，同时绑定在请求上下文
		key := variable.ConfigYml.GetString("Token.BindContextKeyName")
		context.Set(key, loginInfo)

		if err != nil {
			response.ErrorTokenFail(context, consts.PlatformLoginExpiredCode, err.Error())
			return
		}
		//variable.ZapLog.Info("loginInfo:" + loginInfo.Username)

		context.Next()
	}
}
