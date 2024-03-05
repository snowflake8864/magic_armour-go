package authorization

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/response"
	"time"
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

// SaveAgentAccess 保存终端访问信息
func SaveAgentAccess() gin.HandlerFunc {
	return func(context *gin.Context) {
		bodyLogWriter := &responseBody{body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
		context.Writer = bodyLogWriter
		//开始时间
		startTime := time.Now()
		//上下文
		context.Next()

		go func(c *gin.Context) {
			// 结束时间
			endTime := time.Now()

			//返回参数
			responseBody := bodyLogWriter.body.String()

			variable.ZapLog.Info("=============================================SaveAgentAccess start =============================================")
			variable.ZapLog.Info(fmt.Sprintf("请求时间：%s", startTime.Local().Format(variable.DateFormat)))
			variable.ZapLog.Info(fmt.Sprintf("请求地址：%s", c.Request.RequestURI))
			variable.ZapLog.Info(fmt.Sprintf("请求方式：%s", c.Request.Method))
			variable.ZapLog.Info(fmt.Sprintf("请求IP：%s", c.ClientIP()))
			variable.ZapLog.Info(fmt.Sprintf("请求参数：%s", GetRequestParams(c)))
			variable.ZapLog.Info(fmt.Sprintf("响应状态：%d", bodyLogWriter.ResponseWriter.Status()))
			variable.ZapLog.Info(fmt.Sprintf("响应数据：%s", responseBody))
			variable.ZapLog.Info(fmt.Sprintf("请求耗时：%d ms", endTime.Sub(startTime).Milliseconds()))
			variable.ZapLog.Info("=============================================SaveAgentAccess end =============================================")
		}(context.Copy())

	}
}
