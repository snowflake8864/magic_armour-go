package authorization

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_agent"
	access_log_service "goskeleton/app/service/agent/access_log"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/response"
	"net/http"
	"strings"
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
		// token验证通过，同时绑定在请求上下文
		key := variable.ConfigYml.GetString("Token.BindContextKeyName")
		context.Set(key, loginInfo)

		if err != nil {
			response.ErrorTokenFail(context, consts.AgentLoginExpiredCode, err.Error())
			return
		}
		//variable.ZapLog.Info("loginInfo:" + loginInfo.AgentId)

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
			var param = model_agent.AccessLogModel{}
			param.CreatedAt = startTime.Local().Format(variable.DateFormat)
			param.Path = c.Request.RequestURI
			param.NetIp = c.ClientIP()
			param.Method = c.Request.Method
			param.Input = GetRequestParams(c)
			param.HttpCode = bodyLogWriter.ResponseWriter.Status()
			param.DelayTime = time.Now().Sub(startTime).Milliseconds()
			param.Resp = bodyLogWriter.body.String()
			param.Status = enums.IsNot
			if param.HttpCode == http.StatusOK {
				param.Status = enums.IsYes
			}

			//忽略验证的请求地址
			ignoreAuth := map[string]struct{}{
				"/v1/auth": {},
			}
			if _, exist := ignoreAuth[param.Path]; exist {
				param.AgentId = strings.ToLower(c.GetString(consts.ValidatorPrefix + "uid"))
			} else {
				loginInfo := cur_userinfo.GetCurrentAgent(c)
				param.AgentId = loginInfo.AgentId
			}

			if !helpers.IsEmptyString(param.AgentId) {
				access_log_service.CreateCurdFactory().Add(&param)
			}
			/*variable.ZapLog.Info("=============================================SaveAgentAccess start =============================================")
			variable.ZapLog.Info(helpers.JsonEncode(param))
			variable.ZapLog.Info("=============================================SaveAgentAccess end =============================================")*/
		}(context.Copy())

	}
}
