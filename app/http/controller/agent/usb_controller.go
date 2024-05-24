package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/response"
)

// UsbController 外设相关
type UsbController struct {
	ControllerBase
}

// Put 添加设备
func (u *UsbController) Put(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeUsbUpload}
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "data")
	err := agent_post_service.CreateUsbFactory().Save(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// GetBlack 获取禁止列表
func (u *UsbController) GetBlack(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "获取禁止列表")
}

// GetWhite 获取信任列表
func (u *UsbController) GetWhite(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "获取信任列表")
}
