package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/utils/response"
)

// SyslogController 通信日志相关
type SyslogController struct {
}

// Put 上传通讯日志
func (u *SyslogController) Put(ctx *gin.Context) {

	response.SuccessReturnData(ctx, "上传通讯日志")
}
