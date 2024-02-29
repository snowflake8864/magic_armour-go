package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// PortController 端口相关
type PortController struct {
}

// PutAll 卸载回调
func (u *PortController) GetVirtual(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *PortController) PutWorking(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}
