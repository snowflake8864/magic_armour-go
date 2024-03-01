package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// TamperController 防篡改相关
type TamperController struct {
}

// GetDir 获取查询目录
func (u *TamperController) GetDir(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "获取查询目录")
}

// PutDir 上传目录信息
func (u *TamperController) PutDir(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传目录信息")
}

// GetPolicy 下载防篡改策略
func (u *TamperController) GetPolicy(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载防篡改策略")
}
