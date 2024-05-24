package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// ProcSampleController 进程或模块样本相关
type ProcSampleController struct {
	ControllerBase
}

// Get 获取信息
func (u *ProcSampleController) Get(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "获取信息")
}

// Upload 上传样本
func (u *ProcSampleController) Upload(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传")
}
