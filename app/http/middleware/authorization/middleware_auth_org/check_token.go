package middleware_auth_org

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums/enum_org_operate"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

type Params struct {
	Authorization string `header:"Authorization" binding:"required,min=32,max=64"`
}

// CheckToken 检查企业token完整性、有效性中间件
func CheckToken() gin.HandlerFunc {
	//variable.ZapLog.Info("检查终端token完整性、有效性中间件----------------")
	return func(context *gin.Context) {
		headerParams := Params{}
		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.ErrorTokenFail(context, my_errors.OrgLoginExpiredCode, my_errors.JwtTokenMustValid+err.Error())
			return
		}

		isDelay := true
		//忽略延时的请求地址
		/*ignoreDelay := map[string]struct{}{
			"/v5/organization/order/public/expiredRemind": {},
		}
		if _, exist := ignoreDelay[context.Request.RequestURI]; exist {
			isDelay = false
		}*/
		if module, _ := enum_org_operate.GetModuleAndLabelByPath(context.Request.RequestURI); string_helper.IsEmpty(module) {
			isDelay = false
		}

		//variable.ZapLog.Info("Authorization:" + headerParams.Authorization)
		loginInfo, err := user_service.CreateAuthFactory().DecryptByLogin(headerParams.Authorization, isDelay)
		// token验证通过，同时绑定在请求上下文
		key := variable.ConfigYml.GetString("Token.BindContextKeyName")
		context.Set(key, loginInfo)

		if err != nil {
			response.ErrorTokenFail(context, my_errors.OrgLoginExpiredCode, err.Error())
			return
		}
		//variable.ZapLog.Info("loginInfo:" + loginInfo.Username)

		context.Next()
	}
}
