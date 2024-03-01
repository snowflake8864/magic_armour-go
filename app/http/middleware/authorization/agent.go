package authorization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/response"
)

type AgentParams struct {
	Authorization string `header:"Authorization" binding:"required,min=32,max=64"`
}

// CheckAgentToken 检查终端token完整性、有效性中间件
func CheckAgentToken() gin.HandlerFunc {
	//variable.ZapLog.Info("检查终端token完整性、有效性中间件----------------")
	return func(context *gin.Context) {
		headerParams := AgentParams{}
		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.ErrorTokenFail(context, consts.AgentLoginExpiredCode, consts.JwtTokenMustValid+err.Error())
			return
		}
		//variable.ZapLog.Info("Authorization:" + headerParams.Authorization)
		loginInfo, err := agent_service.CreateCacheFactory().DecryptByLogin(headerParams.Authorization, true)
		if err != nil {
			response.ErrorTokenFail(context, consts.AgentLoginExpiredCode, err.Error())
			return
		}
		//variable.ZapLog.Info("loginInfo:" + loginInfo.AgentId)

		key := variable.ConfigYml.GetString("Token.BindContextKeyName")
		// token验证通过，同时绑定在请求上下文
		context.Set(key, loginInfo)

		context.Next()
	}
}
