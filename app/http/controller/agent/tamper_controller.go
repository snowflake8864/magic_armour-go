package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/response"
)

// TamperController 防篡改相关
type TamperController struct {
	ControllerBase
}

// GetDir 获取查询目录
func (u *TamperController) GetDir(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "获取查询目录")
}

// PutDir 上传目录信息
func (u *TamperController) PutDir(ctx *gin.Context) {
	loginInfo := u.GetLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeDirStructureUpload}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "dir")
	err := agent_post_service.CreateAgentFactory().SaveDir(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// GetPolicy 下载防篡改策略
func (u *TamperController) GetPolicy(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载防篡改策略")
}
