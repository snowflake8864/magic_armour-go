package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/response"
	"strings"
)

// AuthController 安全验证相关
type AuthController struct {
}

// Login 登录
func (u *AuthController) Login(ctx *gin.Context) {
	var param = model_agent.AgentModel{}

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
	param.AgentStartAt = helpers.TimestampToDatetime(int64(ctx.GetFloat64(consts.ValidatorPrefixAgent + "asstarttime")))
	param.OsStartAt = helpers.TimestampToDatetime(int64(ctx.GetFloat64(consts.ValidatorPrefixAgent + "osstarttime")))
	param.Auth = ctx.GetString(consts.ValidatorPrefixAgent + "auth")
	param.AdminId = int(ctx.GetFloat64(consts.ValidatorPrefixAgent + "userid"))
	param.NetIp = ctx.ClientIP()

	if helpers.StringIsEmpty(param.AgentIp) {
		param.AgentIp = param.NetIp
	}

	//兼容处理，参考bug 6181
	if helpers.StringIsEmpty(param.OsName) {
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

	isOk, data := agent_service.CreateCurdFactory().RegisterOrLogin(&param)
	if !isOk {
		response.FailSimplify(ctx, data)
	}

	response.SuccessReturnData(ctx, gin.H{"token": data})
}
