package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/response"
)

// DllController 模块相关
type DllController struct {
	ControllerBase
}

// PutAll 上传全部模块
func (u *DllController) PutAll(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.GetLogin(ctx)
	tasks := []int{enum_task.TypeDllSuffixIsDll}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "list")
	err := agent_post_service.CreateProcFactory().ToQuery(enum_process.SourceDll, agent_post_service.TaskAll,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PutRunning 上传运行中模块
func (u *DllController) PutRunning(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.GetLogin(ctx)
	tasks := []int{enum_task.TypeDllRefresh}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "proclist")
	err := agent_post_service.CreateProcFactory().ToQuery(enum_process.SourceDll, agent_post_service.TaskRefresh,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PutAuto 主动上传模块
func (u *DllController) PutAuto(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.GetLogin(ctx)

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "proclist")
	err := agent_post_service.CreateProcFactory().ToQuery(enum_process.SourceDll, agent_post_service.TaskAuto,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// GetWhite 下载模块白名单
func (u *DllController) GetWhite(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载模块白名单")
}

// GetBlack 下载模块黑名单
func (u *DllController) GetBlack(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载模块黑名单")
}
