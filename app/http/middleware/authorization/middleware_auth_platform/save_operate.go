package middleware_auth_platform

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_platform_operate"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/model/platform/operate_models"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/service/platform/platform_operate_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"net/http"
	"time"
)

// SaveOperate 保存平台操作记录
func SaveOperate() gin.HandlerFunc {
	return func(context *gin.Context) {
		bodyLogWriter := &authorization.ResponseBody{Body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
		context.Writer = bodyLogWriter
		//开始时间
		startTime := time.Now()
		//上下文
		context.Next()

		go func(c *gin.Context) {
			var param = operate_models.LogModel{}
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

			//是否需要记录
			param.Module, _ = enum_platform_operate.GetModuleAndLabelByPath(param.Path)
			if string_helper.IsEmpty(param.Module) {
				variable.ZapLog.Debug("SaveOperate path error: " + param.Path)
				return
			}

			//忽略验证的请求地址
			ignoreAuth := map[string]struct{}{
				"/v5/platform/auth/checkUser":     {},
				"/v5/platform/auth/getSliderCode": {},
				"/v5/platform/auth/sendSms":       {},
				"/v5/platform/auth/normalLogin":   {},
			}
			if _, exist := ignoreAuth[param.Path]; exist {
				//获取参数中登录名
				param.CreateName, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(c.GetString(consts.ValidatorPrefixPlatform + "name"))
				if string_helper.IsEmpty(param.CreateName) {
					variable.ZapLog.Error("param name is empty, IP: " + param.NetIp)
					return
				}
				//获取管理员
				adminModel := admin_service.CreateCurdFactory().GetItemByUsername(param.CreateName)
				if string_helper.IsEmpty(adminModel.Username) {
					variable.ZapLog.Error("param name is empty,Name:" + param.CreateName + " IP: " + param.NetIp)
					return
				}
				param.CreateId = adminModel.Id
			} else {
				loginInfo := cur_userinfo.GetCurrentAdmin(c)
				param.CreateId = loginInfo.AdminId
				param.CreateName = loginInfo.Username
			}

			if !string_helper.IsEmpty(param.CreateName) {
				go platform_operate_service.CreateCurdFactory().Add(&param)
			}

			/*
				variable.ZapLog.Info("=============================================SaveOperate start =============================================")
				variable.ZapLog.Info(json_helper.Encode(param))
				variable.ZapLog.Info("=============================================SaveOperate end =============================================")
			*/
		}(context.Copy())

	}
}
