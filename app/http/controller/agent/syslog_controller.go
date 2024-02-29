package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// SyslogController 通信日志相关
type SyslogController struct {
}

// Put 卸载回调
func (u *SyslogController) Put(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "")
}
