package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/http/controller/platform"
	"goskeleton/app/utils/response"
)

// ProcessGlobalController 进程全局名单相关
type ProcessGlobalController struct {
	platform.ControllerBase
}

// PageFromProcess 分页
func (u *ProcessGlobalController) PageFromProcess(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageFromProcess")
}

// PageFromDll 分页
func (u *ProcessGlobalController) PageFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "PageFromDll")
}

// AddFromProcess 添加
func (u *ProcessGlobalController) AddFromProcess(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFromProcess")
}

// AddFromDll 添加
func (u *ProcessGlobalController) AddFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFromDll")
}

// ClearFromProcess 清空
func (u *ProcessGlobalController) ClearFromProcess(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ClearFromProcess")
}

// ClearFromDll 清空
func (u *ProcessGlobalController) ClearFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "ClearFromDll")
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

// Edit 编辑
func (u *ProcessGlobalController) Edit(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Edit")
}

// Delete 删除
func (u *ProcessGlobalController) Delete(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Delete")
}
