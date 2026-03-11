package agent

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
	"strings"
)

// AuthController 安全验证相关
type AuthController struct {
}

// MakeUserid 生成用户秘钥
func (u *AuthController) MakeUserid(ctx *gin.Context) {
	orgId := int64(ctx.GetFloat64(consts.ValidatorPrefixAgent + "orgId"))
	userId := int64(ctx.GetFloat64(consts.ValidatorPrefixAgent + "userId"))
	response.SuccessReturnData(ctx, software_service.CreateDownloadFactory().EncryptConfig(orgId, userId, 0))
}

// Login 登录
func (u *AuthController) Login(ctx *gin.Context) {
	var param = agent_models.AgentModel{}

	param.AgentId = strings.ToLower(ctx.GetString(consts.ValidatorPrefixAgent + "uid"))
	param.HostName = ctx.GetString(consts.ValidatorPrefixAgent + "host_name")
	param.MacId = ctx.GetString(consts.ValidatorPrefixAgent + "macid")
	param.AgentIp = ctx.GetString(consts.ValidatorPrefixAgent + "ip")
	param.Ver = ctx.GetString(consts.ValidatorPrefixAgent + "ver")
	param.PlatformType = int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "type"))
	param.OsName = ctx.GetString(consts.ValidatorPrefixAgent + "os")
	param.MemDesc = ctx.GetString(consts.ValidatorPrefixAgent + "memsize")
	param.CpuDesc = ctx.GetString(consts.ValidatorPrefixAgent + "cpu")
	param.HdDesc = ctx.GetString(consts.ValidatorPrefixAgent + "hdsize")
	param.AgentStartAt = date_helper.TimestampToDatetime(string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixAgent + "asstarttime")))
	param.OsStartAt = date_helper.TimestampToDatetime(string_helper.ToInt64(ctx.GetString(consts.ValidatorPrefixAgent + "osstarttime")))
	param.Auth = ctx.GetString(consts.ValidatorPrefixAgent + "auth")
	param.NetIp = ctx.ClientIP()
	userId := ctx.GetString(consts.ValidatorPrefixAgent + "userid")
	param.ModVer = ctx.GetString(consts.ValidatorPrefixAgent + "mod_ver")
	param.ArchType = int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "arch_type"))

	if string_helper.IsEmpty(param.AgentIp) {
		param.AgentIp = param.NetIp
	}
	if param.ArchType == 0 {
		param.ArchType = enum_software.ArchAmd64
	}

	//兼容处理，参考bug 6181
	if string_helper.IsEmpty(param.OsName) {
		if param.PlatformType == enum_software.PlatformTypeWindows {
			param.OsName = "windows"
		} else {
			param.OsName = "centos"
		}
	}

	//系统类型
	param.OsType = enum_software.GetOsTypeByName(param.OsName)
	if param.OsType == 0 {
		response.FailRoutine(ctx, consts.AgentINSTALLCode, "os unknown: "+param.OsName)
		return
	}

	//解密用户信息
	userInfo, err := software_service.CreateDownloadFactory().DecryptConfig(userId)
	if err != nil {
		//@todo 兼容以前离线版本
		if cur_app.GetAppUsageVersion() == variable.AppUsageVersionOffline {
			dictIds := []string{
				enum_dictionary.AgentCompatibleDefaultOrg,
				enum_dictionary.AgentCompatibleDefaultUser,
			}
			mapDict := dictionary_service.CreateCacheFactory().GetMapByIds(dictIds)
			if val, isOk := mapDict[enum_dictionary.AgentCompatibleDefaultOrg]; !isOk || string_helper.IsEmpty(val) {
				response.FailRoutine(ctx, consts.AgentINSTALLCode, err.Error())
				return
			}
			if val, isOk := mapDict[enum_dictionary.AgentCompatibleDefaultUser]; !isOk || string_helper.IsEmpty(val) {
				response.FailRoutine(ctx, consts.AgentINSTALLCode, err.Error())
				return
			}
			userInfo.OrgId = string_helper.ToInt64(mapDict[enum_dictionary.AgentCompatibleDefaultOrg])
			userInfo.UserId = string_helper.ToInt64(mapDict[enum_dictionary.AgentCompatibleDefaultUser])
			if userInfo.OrgId == 0 || userInfo.UserId == 0 {
				response.FailRoutine(ctx, consts.AgentINSTALLCode, err.Error())
				return
			}
		} else {
			response.FailRoutine(ctx, consts.AgentINSTALLCode, err.Error())
			return
		}
	} else {
		//验证用户
		userModel := user_service.CreateCurdFactory().GetItemByOrgAndUser(userInfo.OrgId, userInfo.UserId)
		if userModel.Id == 0 {
			response.FailRoutine(ctx, consts.AgentINSTALLCode, fmt.Sprintf("user not exist: orgId=%d, userId=%d", userInfo.OrgId, userInfo.UserId))
			return
		}
	}

	//注册登录
	param.OrgId = userInfo.OrgId
	param.UserId = userInfo.UserId
	isOk, data := agent_service.CreateOnlineFactory().RegisterOrLogin(&param)
	if !isOk {
		response.FailRoutine(ctx, consts.AgentParamErrorCode, data)
		return
	}

	response.SuccessReturnData(ctx, gin.H{"token": data})
}
