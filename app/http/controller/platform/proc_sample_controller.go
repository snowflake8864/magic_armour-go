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

// Delete 删除
func (u *ProcSampleController) Delete(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Delete")
}

// Clear 清空
func (u *ProcSampleController) Clear(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Clear")
}

// Download 下载
func (u *ProcSampleController) Download(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Download")
}

// AddFromProcess 添加来源进程
func (u *ProcSampleController) AddFromProcess(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFromProcess")
}

// AddFromDll 添加来源模块
func (u *ProcSampleController) AddFromDll(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "AddFromDll")
}
