package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// ProcessController 进程相关
type ProcessController struct {
	ControllerBase
}

// PutAll 上传全部进程
func (u *ProcessController) PutAll(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传全部进程")
}

// PutRunning 上传运行中进程
func (u *ProcessController) PutRunning(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传运行中进程")
}

// PutAuto 主动上传进程
func (u *ProcessController) PutAuto(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "主动上传进程")
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
