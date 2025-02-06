package middleware_auth_org

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_org_operate"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/middleware/authorization"
	"goskeleton/app/model/organization/org_operate_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/service/organization/org_operate_service"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/cur_userinfo"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/response"
	"time"
)

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
			var param = org_operate_models.LogModel{}
			param.CreatedAt = startTime.Local().Format(variable.DateFormatNow)
			param.Path = c.Request.RequestURI
			param.NetIp = c.ClientIP()
			param.Method = c.Request.Method
			//是否需要记录
			param.Module, _ = enum_org_operate.GetModuleAndLabelByPath(param.Path)
			if string_helper.IsEmpty(param.Module) {
				//variable.ZapLog.Debug("SaveOrgOperate path error: " + param.Path)
				return
			}

			param.Input = authorization.GetRequestParams(c)
			param.HttpCode = bodyLogWriter.ResponseWriter.Status()
			param.DelayTime = time.Now().Sub(startTime).Milliseconds()
			param.Resp = bodyLogWriter.Body.String()
			param.Status = enum_org_operate.StatusFail
			if response.CheckIsSuccess(param.HttpCode, &param.Resp) {
				param.Status = enum_org_operate.StatusSuccess
			}

			//忽略验证的请求地址
			ignoreAuth := map[string]struct{}{
				"/v5/organization/auth/checkUser":     {},
				"/v5/organization/auth/getSliderCode": {},
				"/v5/organization/auth/sendSms":       {},
				"/v5/organization/auth/normalLogin":   {},
			}
			if _, exist := ignoreAuth[param.Path]; exist {
				//获取参数中登录名
				param.CreateName, _ = openssl_encrypt.CreateRsaFactory().DecryptByBase64(c.GetString(consts.ValidatorPrefixOrg + "name"))
				if string_helper.IsEmpty(param.CreateName) {
					variable.ZapLog.Error("SaveOrgOperate param name is empty, IP: " + param.NetIp)
					return
				}
				//获取用户
				userModel := user_service.CreateCurdFactory().GetItemByName(param.CreateName)
				if string_helper.IsEmpty(userModel.Username) {
					variable.ZapLog.Error("SaveOrgOperate param name is empty,Name:" + param.CreateName + " ,IP: " + param.NetIp)
					return
				}
				param.CreateId = userModel.Id
				param.OrgId = userModel.OrgId
				param.Source = enum_org_operate.SourceWeb
			} else if param.Path == "/v5/organization/auth/hsLogin" { //幻势秘钥登录
				//获取参数中企业标识
				orgCode := huanshi_platform.CreateCallbackFactory().GetDeviceIdByToken(c.GetString(consts.ValidatorPrefixOrg + "token"))
				if string_helper.IsEmpty(orgCode) {
					variable.ZapLog.Error("SaveOrgOperate param token is empty, IP: " + param.NetIp)
					return
				}
				orgModel, err := organization_models.CreateOrganizationFactory().GetItemByCode(orgCode)
				if err != nil || orgModel.Id == 0 {
					variable.ZapLog.Error("SaveOrgOperate param name is empty,orgCode:" + orgCode + " ,IP: " + param.NetIp)
					return
				}
				param.CreateId = 0
				param.CreateName = orgModel.UserName
				param.OrgId = orgModel.Id
				param.Source = enum_org_operate.SourceHs
			} else {
				if loginInfo := cur_userinfo.GetCurrentUser(c); loginInfo != nil && loginInfo.OrgId > 0 {
					param.OrgId = loginInfo.OrgId
					param.CreateId = loginInfo.UserId
					param.CreateName = loginInfo.Username
					param.Source = loginInfo.Source
				}
			}

			if param.OrgId > 0 && !string_helper.IsEmpty(param.CreateName) {
				go org_operate_service.CreateCurdFactory().Add(&param)
			}
			/*
				variable.ZapLog.Info("=============================================SaveOrgOperate start =============================================")
				variable.ZapLog.Info(json_helper.Encode(param))
				variable.ZapLog.Info("=============================================SaveOrgOperate end =============================================")

			*/

		}(context.Copy())

	}
}
