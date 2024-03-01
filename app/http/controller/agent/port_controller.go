package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// PortController 端口相关
type PortController struct {
}

// GetVirtual 下载虚开端口
func (u *PortController) GetVirtual(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载虚开端口")
}

// PutWorking 上传业务端口
func (u *PortController) PutWorking(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传业务端口")
}
