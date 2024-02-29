package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// BlackmailController 勒索相关
type BlackmailController struct {
}

// PutAll 卸载回调
func (u *BlackmailController) GetPolicy(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}
