package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_pw_jump"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/utils/response"
)

// JumpController 封堵相关
type JumpController struct {
	ControllerBase
}

// GetPwJump 获取密码跳变
func (u *JumpController) GetPwJump(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypePwJump}, ctx.ClientIP())

	config := enum_pw_jump.PwJumpConfig{}
	//优先终端策略
	if config = agent_service.CreateExtendFactory().GetPwJumpConfig(loginInfo.OrgId, loginInfo.AgentId); config.Switch == enums.SwitchOff {
		OrgExtendFact := organize_service.CreateOrgExtendFactory()
		config = OrgExtendFact.FormatPwJumpConfig(OrgExtendFact.GetItemByOrg(loginInfo.OrgId).IpJumpConfig)
		if config.Switch == enums.SwitchOff {
			response.FailSimplify(ctx, "Switch Off")
			return
		}
	}

	pw := config.GeneratePw()
	//保存密码
	if err := agent_service.CreateExtendFactory().SetPwJumpNewest(loginInfo.OrgId, loginInfo.AgentId, pw); err != nil {
		response.FailSimplify(ctx, "save pw fail")
		return
	}

	response.SuccessReturnData(ctx, gin.H{
		"pw": pw,
	})
}

// GetIpJump 获取IP跳变
func (u *JumpController) GetIpJump(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeIpJump}, ctx.ClientIP())

	response.SuccessReturnData(ctx, gin.H{
		"source_ip": ctx.ClientIP(),
		"target_ip": ctx.ClientIP(),
	})
}
