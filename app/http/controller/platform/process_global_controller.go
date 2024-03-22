package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// ProcessGlobalController 进程全局名单相关
type ProcessGlobalController struct {
	ControllerBase
}

// Page 分页
func (u *ProcessGlobalController) Page(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Page")
}

// BatchWhiteFromProcess 批量白来源进程
func (u *ProcessGlobalController) BatchWhiteFromProcess(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchWhiteFromProcess")
}

// BatchBlackFromProcess 批量黑来源进程
func (u *ProcessGlobalController) BatchBlackFromProcess(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchBlackFromProcess")
}

// BatchWhiteFromDll 批量白来源模块
func (u *ProcessGlobalController) BatchWhiteFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchWhiteFromDll")
}

// BatchBlackFromDll 批量黑来源模块
func (u *ProcessGlobalController) BatchBlackFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "BatchBlackFromDll")
}
