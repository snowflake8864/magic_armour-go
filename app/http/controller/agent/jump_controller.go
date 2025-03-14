package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/agent/task_service"
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

	pw := "d2x2d2w22s24#$@GGG3"
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
