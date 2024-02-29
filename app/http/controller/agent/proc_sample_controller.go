package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// ProcSampleController 进程或模块样本相关
type ProcSampleController struct {
}

// PutAll 卸载回调
func (u *ProcSampleController) Get(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *ProcSampleController) Upload(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}
