package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// BlackmailController 勒索相关
type BlackmailController struct {
	ControllerBase
}

// GetPolicy 下载勒索保护策略
func (u *BlackmailController) GetPolicy(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "下载勒索保护策略")
}
