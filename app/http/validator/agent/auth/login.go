package auth

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/http/controller/agent"
	"goskeleton/app/http/validator/core/data_transfer"
	"goskeleton/app/utils/helpers"
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
}

func (n Login) CheckParams(context *gin.Context) {
	//1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&n); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	if helpers.IsEmptyString(n.Uuid) {
		response.FailRoutine(context, consts.AgentINSTALLCode, "param [Uuid] not empty")
		return
	}
	if helpers.IsEmptyString(n.MacId) {
		response.FailRoutine(context, consts.AgentINSTALLCode, "param [MacId] not empty")
		return
	}
	if helpers.IsEmptyString(n.Ver) {
		response.FailRoutine(context, consts.AgentINSTALLCode, "param [Ver] not empty")
		return
	}
	if helpers.IsEmptyString(n.Auth) {
		response.FailRoutine(context, consts.AgentINSTALLCode, "param [Auth] error")
		return
	}
	if helpers.IsEmptyString(n.UserId) {
		response.FailRoutine(context, consts.AgentINSTALLCode, "param [UserId] error")
		return
	}

	if label := enums.GET1stIsStrAnd2ndIsIntReturnStr("platform_type", n.PlatformType, enum_software.EnumMap); helpers.IsEmptyString(label) {
		response.FailRoutine(context, consts.AgentINSTALLCode, "param [PlatformType] error")
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(n, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.FailSimplify(context, "agent auth logon表单验证器json化失败")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		(&agent.AuthController{}).Login(extraAddBindDataContext)
	}

}
