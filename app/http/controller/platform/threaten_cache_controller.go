package platform

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// ThreatenCacheController 威胁情报缓存相关
type ThreatenCacheController struct {
	ControllerBase
}

// Levels 等级
func (u *ThreatenCacheController) Levels(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Levels")
}

// Page 分页
func (u *ThreatenCacheController) Page(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "Page")
}
