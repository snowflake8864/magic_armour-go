package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// DllController 模块相关
type DllController struct {
}

// PutAll 上传全部模块
func (u *DllController) PutAll(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传全部模块")
}

// PutRunning 上传运行中模块
func (u *DllController) PutRunning(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传运行中模块")
}

// PutAuto 主动上传模块
func (u *DllController) PutAuto(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "主动上传模块")
}

// GetWhite 下载模块白名单
func (u *DllController) GetWhite(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载模块白名单")
}

// GetBlack 下载模块黑名单
func (u *DllController) GetBlack(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载模块黑名单")
}
