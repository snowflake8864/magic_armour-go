package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// PlugController 封堵相关
type PlugController struct {
}

// GetDynamic 下载动态封堵
func (u *PlugController) GetDynamic(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutDynamic 上传动态封堵(调试用)
func (u *PlugController) PutDynamic(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传动态封堵(调试用)")
}

// GetIpWhite 下载封堵ip白名单(不在使用)
func (u *PlugController) GetIpWhite(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载封堵ip白名单(不在使用)")
}

// GetIpBlack 下载封堵ip黑名单
func (u *PlugController) GetIpBlack(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载封堵ip黑名单")
}
