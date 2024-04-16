package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/platform"
	"goskeleton/app/utils/response"
)

// ThreatenPushController 威胁情报推送相关
type ThreatenPushController struct {
	platform.ControllerBase
}

// Page 分页
func (u *ThreatenPushController) Page(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Page")
}

// AddFromProcess 添加来源进程
func (u *ThreatenPushController) AddFromProcess(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFromProcess")
}

// AddFromDll 添加来源模块
func (u *ThreatenPushController) AddFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFromDll")
}

// BatchWhiteFromProcess 批量白来源进程
func (u *ThreatenPushController) BatchWhiteFromProcess(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchWhiteFromProcess")
}

// BatchBlackFromProcess 批量黑来源进程
func (u *ThreatenPushController) BatchBlackFromProcess(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchBlackFromProcess")
}

// BatchWhiteFromDll 批量白来源模块
func (u *ThreatenPushController) BatchWhiteFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchWhiteFromDll")
}

// BatchBlackFromDll 批量黑来源模块
func (u *ThreatenPushController) BatchBlackFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchBlackFromDll")
}
