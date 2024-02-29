package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// TamperController 防篡改相关
type TamperController struct {
}

// PutAll 卸载回调
func (u *TamperController) GetDir(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *TamperController) PutDir(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutAll 卸载回调
func (u *TamperController) GetPolicy(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}
