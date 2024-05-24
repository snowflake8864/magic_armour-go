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

// ProcessController 进程相关
type ProcessController struct {
	ControllerBase
}

// PutAll 上传全部进程
func (u *ProcessController) PutAll(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	tasks := []int{enum_task.TypeProcessSuffixIsExe}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "list")
	err := agent_post_service.CreateProcFactory().ToQuery(enum_process.SourceProcess, agent_post_service.ProcTaskAll,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PutRunning 上传运行中进程
func (u *ProcessController) PutRunning(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	tasks := []int{enum_task.TypeProcessRefresh}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "proclist")
	err := agent_post_service.CreateProcFactory().ToQuery(enum_process.SourceProcess, agent_post_service.ProcTaskRefresh,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PutAuto 主动上传进程
func (u *ProcessController) PutAuto(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "proclist")
	err := agent_post_service.CreateProcFactory().ToQuery(enum_process.SourceProcess, agent_post_service.ProcTaskAuto,
		loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// GetWhite 下载进程白名单
func (u *ProcessController) GetWhite(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载进程白名单")
}

// GetBlack 下载进程黑名单
func (u *ProcessController) GetBlack(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载进程黑名单")
}

// GetAssign 指定进程:下载信息
func (u *ProcessController) GetAssign(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "定进程:下载信息")
}

// PutAssign 指定进程:上传对应模块
func (u *ProcessController) PutAssign(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "指定进程:上传对应模块")
}
