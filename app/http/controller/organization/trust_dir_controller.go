package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// TrustDirController 信任目录相关
type TrustDirController struct {
	ControllerBase
}

// PageForGlobal 全局： 分页
func (u *TrustDirController) PageForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FilePowers")
}

// AddForGlobal 全局： 添加
func (u *TrustDirController) AddForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FilePowers")
}

// EditForGlobal 全局： 编辑
func (u *TrustDirController) EditForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FilePowers")
}

// SetSwitchForGlobal 全局： 设置开关
func (u *TrustDirController) SetSwitchForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FilePowers")
}

// DeleteForGlobal 全局： 删除
func (u *TrustDirController) DeleteForGlobal(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "FilePowers")
}
