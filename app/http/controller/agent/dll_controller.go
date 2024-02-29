package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// DllController 模块相关
type DllController struct {
}

// PutAll 卸载回调
func (u *DllController) PutAll(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *DllController) PutRunning(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *DllController) PutAuto(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// GetWhite 卸载回调
func (u *DllController) GetWhite(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *DllController) GetBlack(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}
