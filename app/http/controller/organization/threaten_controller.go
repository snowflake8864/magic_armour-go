package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// ThreatenController 威胁情报相关
type ThreatenController struct {
	ControllerBase
}

// GetLevels 通过hashes获取集合
func (u *ThreatenController) GetLevels(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetLevels")
}

// GetAttr 获取属性
func (u *ThreatenController) GetAttr(ctx *gin.Context) {
	response.SuccessReturnData(ctx, "GetLevels")
}
