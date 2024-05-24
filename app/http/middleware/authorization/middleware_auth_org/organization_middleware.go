package middleware_auth_org

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_org_operate"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/model/organization/org_operate_models"
	"goskeleton/app/service/organization/org_operate_service"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
	"net/http"
	"time"
)

type OrganizationParams struct {
	Authorization string `header:"Authorization" binding:"required,min=32,max=64"`
}

// CheckToken 检查企业token完整性、有效性中间件
func CheckToken() gin.HandlerFunc {
	//variable.ZapLog.Info("检查终端token完整性、有效性中间件----------------")
	return func(context *gin.Context) {
		headerParams := OrganizationParams{}
		//  推荐使用 ShouldBindHeader 方式获取头参数
		if err := context.ShouldBindHeader(&headerParams); err != nil {
			response.ErrorTokenFail(context, consts.OrgLoginExpiredCode, consts.JwtTokenMustValid+err.Error())
			return
		}
		//variable.ZapLog.Info("Authorization:" + headerParams.Authorization)
		loginInfo, err := user_service.CreateAuthFactory().DecryptByLogin(headerParams.Authorization, true)
		// token验证通过，同时绑定在请求上下文
		key := variable.ConfigYml.GetString("Token.BindContextKeyName")
		context.Set(key, loginInfo)

		if err != nil {
			response.ErrorTokenFail(context, consts.OrgLoginExpiredCode, err.Error())
			return
		}
		//variable.ZapLog.Info("loginInfo:" + loginInfo.Username)

		context.Next()
	}
}

// SaveOperate 保存企业操作记录
func SaveOperate() gin.HandlerFunc {
	return func(context *gin.Context) {
		bodyLogWriter := &authorization.ResponseBody{Body: bytes.NewBufferString(""), ResponseWriter: context.Writer}
		context.Writer = bodyLogWriter

		//上下文
		context.Next()
		//开始时间
		startTime := time.Now()
		go func(c *gin.Context) {
			var param = org_operate_models.OperateLogModel{}
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
			param.Module, _ = enum_org_operate.GetModuleAndLabelByPath(param.Path)
			if string_helper.IsEmpty(param.Module) {
				variable.ZapLog.Debug("SaveOrgOperate path error: " + param.Path)
				return
			}

			//忽略验证的请求地址
			ignoreAuth := map[string]struct{}{
				"/organization/auth/checkUser":     {},
				"/organization/auth/getSliderCode": {},
				"/organization/auth/sendSms":       {},
				"/organization/auth/normalLogin":   {},
				"/organization/auth/logout":        {},
				"/organization/auth/hsLogin":       {},
			}
			if _, exist := ignoreAuth[param.Path]; exist {
				//获取参数中登录名
				param.CreateName, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(c.GetString(consts.ValidatorPrefixPlatform + "name"))
				if string_helper.IsEmpty(param.CreateName) {
					variable.ZapLog.Error("SaveOrgOperate param name is empty, IP: " + param.NetIp)
					return
				}
				//获取用户
				userModel := user_service.CreateCurdFactory().GetItemByName(param.CreateName)
				if string_helper.IsEmpty(userModel.Username) {
					variable.ZapLog.Error("SaveOrgOperate param name is empty,Name:" + param.CreateName + " IP: " + param.NetIp)
					return
				}
				param.CreateId = userModel.Id
				param.OrgId = userModel.OrgId
			} else {
				loginInfo := cur_userinfo.GetCurrentUser(c)
				param.OrgId = loginInfo.OrgId
				param.CreateId = loginInfo.UserId
				param.CreateName = loginInfo.Username
			}

			if !string_helper.IsEmpty(param.CreateName) {
				org_operate_service.CreateCurdFactory().Add(&param)
			}

			/*
				variable.ZapLog.Info("=============================================SaveOrgOperate start =============================================")
				variable.ZapLog.Info(json_helper.Encode(param))
				variable.ZapLog.Info("=============================================SaveOrgOperate end =============================================")
			*/
		}(context.Copy())

	}
}
