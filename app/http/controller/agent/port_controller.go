package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/response"
)

// PortController 端口相关
type PortController struct {
	ControllerBase
}

// GetVirtual 下载虚开端口
func (u *PortController) GetVirtual(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "获取虚开规则")
}

// PutWorking 上传业务端口
func (u *PortController) PutWorking(ctx *gin.Context) {
	loginInfo := u.GetLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypePortWorkingUpload}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "service")
	err := agent_post_service.CreatePortFactory().SaveWorking(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)

}
