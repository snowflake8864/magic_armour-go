package authorization

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/response"
	"net/http"
	"time"
)

type PlatformParams struct {
	Authorization string `header:"Authorization" binding:"required,min=32,max=64"`
}

// CheckPlatformPToken 检查平台token完整性、有效性中间件
func CheckPlatformPToken() gin.HandlerFunc {
	//variable.ZapLog.Info("检查终端token完整性、有效性中间件----------------")
	return func(context *gin.Context) {
		headerParams := AgentParams{}
		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.ErrorTokenFail(context, consts.AgentLoginExpiredCode, consts.JwtTokenMustValid+err.Error())
			return
		}
		//variable.ZapLog.Info("Authorization:" + headerParams.Authorization)
		loginInfo, err := admin_service.CreateAuthFactory().DecryptByLogin(headerParams.Authorization, true)
		// token验证通过，同时绑定在请求上下文
		key := variable.ConfigYml.GetString("Token.BindContextKeyName")
		context.Set(key, loginInfo)

		if err != nil {
			response.ErrorTokenFail(context, consts.AgentLoginExpiredCode, err.Error())
			return
		}
		variable.ZapLog.Info("loginInfo:" + loginInfo.Username)

		context.Next()
	}
}

// SavePlatformOperate 保存平台操作记录
func SavePlatformOperate() gin.HandlerFunc {
	return func(context *gin.Context) {
		bodyLogWriter := &responseBody{body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
		context.Writer = bodyLogWriter
		//开始时间
		startTime := time.Now()
		//上下文
		context.Next()

		go func(c *gin.Context) {
			var param = model_platform.OperateLogModel{}
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
				"/platform/auth/checkUser":     {},
				"/platform/auth/getSliderCode": {},
				"/platform/auth/sendSms":       {},
				"/platform/auth/normalLogin":   {},
				"/platform/auth/logout":        {},
				"/platform/auth/hsLogin":       {},
			}
			if _, exist := ignoreAuth[param.Path]; exist {

			} else {
				loginInfo := cur_userinfo.GetCurrentAdmin(c)
				param.CreateId = loginInfo.AdminId
				param.CreateName = loginInfo.Username
			}

			variable.ZapLog.Info("=============================================SavePlatformOperate start =============================================")
			variable.ZapLog.Info(helpers.JsonEncode(param))
			variable.ZapLog.Info("=============================================SavePlatformOperate end =============================================")
		}(context.Copy())

	}
}
