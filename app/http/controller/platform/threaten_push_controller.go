package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// ThreatenPushController 威胁情报推送相关
type ThreatenPushController struct {
	ControllerBase
}

// Page 分页
func (u *ThreatenPushController) Page(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Page")
}

// AllWhite 全部白
func (u *ThreatenPushController) AllWhite(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AllWhite")
}

// BatchWhite 批量白
func (u *ThreatenPushController) BatchWhite(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchWhite")
}

// BatchBlack 批量黑
func (u *ThreatenPushController) BatchBlack(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchBlack")
}

// Clear 清空
func (u *ThreatenPushController) Clear(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Clear")
}
