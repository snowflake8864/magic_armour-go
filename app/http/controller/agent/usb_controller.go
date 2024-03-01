package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// UsbController 外设相关
type UsbController struct {
}

// Put 添加设备
func (u *UsbController) Put(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "添加设备")
}

// GetBlack 获取禁止列表
func (u *UsbController) GetBlack(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "获取禁止列表")
}

// GetWhite 获取信任列表
func (u *UsbController) GetWhite(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "获取信任列表")
}
