package auth

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/http/controller/agent"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

type Login struct {
	Uuid         string `form:"uid" json:"uid"`
	HostName     string `form:"host_name" json:"host_name"`
	MacId        string `form:"macid" json:"macid"`
	AgentIp      string `form:"ip" json:"ip"`
	Ver          string `form:"ver" json:"ver"`
	PlatformType int    `form:"type" json:"type"`
	OsName       string `form:"os" json:"os"`
	MemSize      string `form:"memsize" json:"memsize"`
	Cpu          string `form:"cpu" json:"cpu"`
	HdSize       string `form:"hdsize" json:"hdsize"`
	AsStartTime  string `form:"asstarttime" json:"asstarttime"`
	OsStartTime  string `form:"osstarttime" json:"osstarttime"`
	Auth         string `form:"auth" json:"auth"`
	UserId       string `form:"userid" json:"userid"`
	ModVer       string `form:"mod_ver" json:"mod_ver"`
	ArchType     int    `form:"arch_type" json:"arch_type"`
}

func (n Login) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	if string_helper.IsEmpty(n.Uuid) {
		response.FailRoutine(context, my_errors.AgentINSTALLCode, "param [Uuid] not empty")
		return
	}
	if string_helper.IsEmpty(n.MacId) {
		response.FailRoutine(context, my_errors.AgentINSTALLCode, "param [MacId] not empty")
		return
	}
	if string_helper.IsEmpty(n.Ver) {
		response.FailRoutine(context, my_errors.AgentINSTALLCode, "param [Ver] not empty")
		return
	}
	if string_helper.IsEmpty(n.Auth) {
		response.FailRoutine(context, my_errors.AgentINSTALLCode, "param [Auth] error")
		return
	}
	if string_helper.IsEmpty(n.UserId) {
		response.FailRoutine(context, my_errors.AgentINSTALLCode, "param [UserId] error")
		return
	}

	if label := enums.Get1stIsStrAnd2ndIsIntReturnStr("platform_type", n.PlatformType, enum_software.EnumMap); string_helper.IsEmpty(label) {
		response.FailRoutine(context, my_errors.AgentINSTALLCode, "param [PlatformType] error")
		return
	}
	if n.ArchType > 0 {
		if label := enums.Get1stIsStrAnd2ndIsIntReturnStr("arch_type", n.ArchType, enum_software.EnumMap); string_helper.IsEmpty(label) {
			response.FailRoutine(context, my_errors.AgentINSTALLCode, "param [ArchType] error")
			return
		}
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefixAgent, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "[agent auth logon]表单验证器json化失败")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&agent.AuthController{}).Login(extraAddBindDataContext)
	}

}
