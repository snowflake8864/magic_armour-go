package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// AlarmController 告警相关
type AlarmController struct {
}

// PutMaster 卸载回调
func (u *AlarmController) PutMaster(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}

// PutPort 卸载回调
func (u *AlarmController) PutPort(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}
