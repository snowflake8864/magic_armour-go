package middleware_auth_agent

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/model/agent/agent_access_models"
	accessLogService "goskeleton/app/service/agent/access_log_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/string_helper"
	"net/http"
	"strings"
	"time"
)

// SaveAccess 保存终端访问信息
func SaveAccess() gin.HandlerFunc {
	return func(context *gin.Context) {
		bodyLogWriter := &authorization.ResponseBody{Body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
		context.Writer = bodyLogWriter
		//开始时间
		startTime := time.Now()
		//上下文
		context.Next()

		go func(c *gin.Context) {
			var param = agent_access_models.AccessLogModel{}
			param.CreatedAt = startTime.Local().Format(variable.DateFormat)
			param.Path = c.Request.RequestURI
			param.NetIp = c.ClientIP()
			param.Method = c.Request.Method
			param.Input = authorization.GetRequestParams(c)
			param.HttpCode = bodyLogWriter.ResponseWriter.Status()
			param.DelayTime = time.Now().Sub(startTime).Milliseconds()
			param.Resp = bodyLogWriter.Body.String()
			param.Status = enums.IsNot
			if param.HttpCode == http.StatusOK {
				param.Status = enums.IsYes
			}

			//忽略验证的请求地址
			ignoreAuth := map[string]struct{}{
				"/v1/auth":      {},
				"/v1/getuserid": {},
			}
			if _, exist := ignoreAuth[param.Path]; exist {
				param.AgentId = strings.ToLower(c.GetString(consts.ValidatorPrefixAgent + "uid"))
			} else {
				loginInfo := cur_userinfo.GetCurrentAgent(c)
				param.AgentId = loginInfo.AgentId
			}

			if !string_helper.IsEmpty(param.AgentId) {
				accessLogService.CreateCurdFactory().Add(&param)
			}
			/*
				variable.ZapLog.Info("=============================================SaveAccess start =============================================")
				variable.ZapLog.Info(json_helper.Encode(param))
				variable.ZapLog.Info("=============================================SaveAccess end =============================================")
			*/
		}(context.Copy())

	}
}
