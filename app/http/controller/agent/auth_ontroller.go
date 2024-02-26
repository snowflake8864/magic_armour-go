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

// Auth 通过秘钥获取口令
func (u *AuthController) Auth(ctx *gin.Context) {
	var param = model_agent.AgentModel{}

	param.AgentId = strings.ToLower(ctx.GetString(consts.ValidatorPrefix + "uid"))
	param.HostName = ctx.GetString(consts.ValidatorPrefix + "host_name")
	param.HostName = helpers.ConvertToString(param.HostName, "gbk", "utf-8")
	param.HostName = helpers.ConvertToString(param.HostName, "GB2312", "utf-8")
	param.HostName = helpers.ConvertToString(param.HostName, "BIG5", "utf-8")
	param.MacId = ctx.GetString(consts.ValidatorPrefix + "macid")
	param.AgentIp = ctx.GetString(consts.ValidatorPrefix + "ip")
	param.Ver = ctx.GetString(consts.ValidatorPrefix + "ver")
	param.PlatformType = int(ctx.GetFloat64(consts.ValidatorPrefix + "type"))
	param.OsName = ctx.GetString(consts.ValidatorPrefix + "os")
	param.MemSize = ctx.GetString(consts.ValidatorPrefix + "memsize")
	param.CpuSize = ctx.GetString(consts.ValidatorPrefix + "cpu")
	param.HdSize = ctx.GetString(consts.ValidatorPrefix + "hdsize")
	param.AgentStartAt = ctx.GetString(consts.ValidatorPrefix + "asstarttime")
	param.Auth = ctx.GetString(consts.ValidatorPrefix + "auth")
	param.AdminId = int(ctx.GetFloat64(consts.ValidatorPrefix + "userid"))
	param.NetIp = ctx.ClientIP()

	if helpers.IsEmptyString(param.AgentIp) {
		param.AgentIp = param.NetIp
	}

	//兼容处理，参考bug 6181
	if helpers.IsEmptyString(param.OsName) {
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

	isOk, data := agent_service.CreateLoginFactory().RegisterOrLogin(&param)
	if !isOk {
		response.FailSimplify(ctx, data)
	}

	response.SuccessReturnData(ctx, gin.H{"token": data})
}
