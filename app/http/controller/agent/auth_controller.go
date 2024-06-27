package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/service/organization/user_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/software_service"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
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
	param.AgentStartAt = date_helper.TimestampToDatetime(int64(ctx.GetFloat64(consts.ValidatorPrefixAgent + "asstarttime")))
	param.OsStartAt = date_helper.TimestampToDatetime(int64(ctx.GetFloat64(consts.ValidatorPrefixAgent + "osstarttime")))
	param.Auth = ctx.GetString(consts.ValidatorPrefixAgent + "auth")
	param.NetIp = ctx.ClientIP()
	userId := ctx.GetString(consts.ValidatorPrefixAgent + "userid")

	if string_helper.IsEmpty(param.AgentIp) {
		param.AgentIp = param.NetIp
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
		response.FailRoutine(ctx, consts.AgentINSTALLCode, "未知的系统类型")
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
			userInfo.OrgId = string_helper.ToInt64(mapDict[enum_dictionary.AgentCompatibleDefaultOrg])
			userInfo.UserId = string_helper.ToInt64(mapDict[enum_dictionary.AgentCompatibleDefaultUser])
		} else {
			response.FailRoutine(ctx, consts.AgentINSTALLCode, err.Error())
			return
		}
	}
	errOrg := int_helper.Int64ToStr(userInfo.OrgId) + "+" + int_helper.Int64ToStr(userInfo.UserId)
	//验证企业
	orgModel := organize_service.CreateOrgCurdFactory().GetItemByOrg(userInfo.OrgId)
	if orgModel.Id == 0 {
		response.FailRoutine(ctx, consts.AgentINSTALLCode, "企业不存在:"+errOrg)
		return
	}
	if orgModel.Status != enum_organization.StatusNormal {
		response.FailRoutine(ctx, consts.AgentINSTALLCode, "企业状态异常:"+errOrg)
		return
	}
	//验证用户
	userModel := user_service.CreateCurdFactory().GetItemByOrgAndUser(userInfo.OrgId, userInfo.UserId)
	if userModel.Id == 0 {
		response.FailRoutine(ctx, consts.AgentINSTALLCode, "用户不存在:"+errOrg)
		return
	}

	//新增或修改终端
	param.OrgId = userInfo.OrgId
	param.UserId = userInfo.OrgId
	isOk, data := agent_service.CreateOnlineFactory().RegisterOrLogin(&param)
	if !isOk {
		response.FailSimplify(ctx, data)
		return
	}

	response.SuccessReturnData(ctx, gin.H{"token": data})
}
