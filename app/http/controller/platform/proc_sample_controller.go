package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// ProcSampleController 进程或模块样本相关
type ProcSampleController struct {
	ControllerBase
}

// Page 分页
func (u *ProcSampleController) Page(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Page")
}

// Add 添加
func (u *ProcSampleController) Add(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Add")
}

// Edit 编辑
func (u *ProcSampleController) Edit(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Edit")
}

// Delete 删除
func (u *ProcSampleController) Delete(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Delete")
}

// Clear 清空
func (u *ProcSampleController) Clear(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Clear")
}
