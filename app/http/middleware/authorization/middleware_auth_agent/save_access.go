package middleware_auth_agent

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent_access"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/model/agent/agent_access_models"
	"goskeleton/app/service/agent/access_log_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/software_service"
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
			if val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.AgentSaveAccessLog); string_helper.ToInt(val) == enums.SwitchOff {
				return
			}

			var param = agent_access_models.AccessLogModel{}
			param.CreatedAt = startTime.Local().Format(variable.DateFormat)
			param.Path = c.Request.RequestURI
			param.NetIp = c.ClientIP()
			param.Method = c.Request.Method
			param.Input = authorization.GetRequestParams(c)
			param.HttpCode = bodyLogWriter.ResponseWriter.Status()
			param.DelayTime = time.Now().Sub(startTime).Milliseconds()
			param.Resp = bodyLogWriter.Body.String()
			param.Status = enum_agent_access.StatusFail
			if param.HttpCode == http.StatusOK {
				param.Status = enum_agent_access.StatusSuccess
			}

			//忽略验证的请求地址
			ignoreAuth := map[string]struct{}{
				"/v1/auth":      {},
				"/v1/getuserid": {},
			}
			if _, exist := ignoreAuth[param.Path]; exist {
				//解密用户信息
				userId := c.GetString(consts.ValidatorPrefixAgent + "userid")
				if userInfo, err := software_service.CreateDownloadFactory().DecryptConfig(userId); err == nil {
					param.OrgId = userInfo.OrgId
					param.AgentId = strings.ToLower(c.GetString(consts.ValidatorPrefixAgent + "uid"))
				} else {
					return
				}
			} else {
				if loginInfo := cur_userinfo.GetCurrentAgent(c); loginInfo != nil && loginInfo.OrgId > 0 {
					param.OrgId = loginInfo.OrgId
					param.AgentId = loginInfo.AgentId
				}
			}

			if param.OrgId > 0 && !string_helper.IsEmpty(param.AgentId) {
				access_log_service.CreateCurdFactory().Add(&param)
			}
			/*
				variable.ZapLog.Info("=============================================SaveAccess start =============================================")
				variable.ZapLog.Info(json_helper.Encode(param))
				variable.ZapLog.Info("=============================================SaveAccess end =============================================")
			*/
		}(context.Copy())

	}
}
