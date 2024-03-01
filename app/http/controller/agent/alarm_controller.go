package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// AlarmController 告警相关
type AlarmController struct {
}

// PutMaster 上传报警信息(除端口虚开外)
func (u *AlarmController) PutMaster(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传报警信息(除端口虚开外)")
}

// PutPort 上传端口虚开告警
func (u *AlarmController) PutPort(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传端口虚开告警")
}
