package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// ProcessController 进程相关
type ProcessController struct {
}

// PutAll 卸载回调
func (u *ProcessController) PutAll(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *ProcessController) PutRunning(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *ProcessController) PutAuto(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// GetWhite 卸载回调
func (u *ProcessController) GetWhite(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *ProcessController) GetBlack(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *ProcessController) GetAssign(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *ProcessController) PutAssign(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}
