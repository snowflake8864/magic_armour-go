package authorization

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_platform_operate"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/platform/operate_models"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/service/platform/operate_log_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
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
		headerParams := PlatformParams{}
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
			var param = operate_models.OperateLogModel{}
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

			//是否需要记录
			param.Module, _ = enum_platform_operate.GetModuleAndLabelByPath(param.Path)
			if string_helper.IsEmpty(param.Module) {
				variable.ZapLog.Debug("SavePlatformOperate path error: " + param.Path)
				return
			}

			//忽略验证的请求地址
			ignoreAuth := map[string]struct{}{
				"/platform/auth/checkUser":     {},
				"/platform/auth/getSliderCode": {},
				"/platform/auth/sendSms":       {},
				"/platform/auth/normalLogin":   {},
				"/platform/auth/logout":        {},
			}
			if _, exist := ignoreAuth[param.Path]; exist {
				//获取参数中登录名
				param.CreateName, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(c.GetString(consts.ValidatorPrefixPlatform + "name"))
				if string_helper.IsEmpty(param.CreateName) {
					variable.ZapLog.Error("SavePlatformOperate param name is empty, IP: " + param.NetIp)
					return
				}
				//获取管理员
				adminModel := admin_service.CreateCurdFactory().GetItemByUsername(param.CreateName)
				if string_helper.IsEmpty(adminModel.Username) {
					variable.ZapLog.Error("SavePlatformOperate param name is empty,Name:" + param.CreateName + " IP: " + param.NetIp)
					return
				}
				param.CreateId = adminModel.Id
			} else {
				loginInfo := cur_userinfo.GetCurrentAdmin(c)
				param.CreateId = loginInfo.AdminId
				param.CreateName = loginInfo.Username
			}

			if !string_helper.IsEmpty(param.CreateName) {
				operate_log_service.CreateCurdFactory().Add(&param)
			}

			/*variable.ZapLog.Info("=============================================SavePlatformOperate start =============================================")
			variable.ZapLog.Info(helpers.JsonEncode(param))
			variable.ZapLog.Info("=============================================SavePlatformOperate end =============================================")*/
		}(context.Copy())

	}
}
