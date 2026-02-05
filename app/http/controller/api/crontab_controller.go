package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/utils/response"
	"goskeleton/command/console"
)

// CrontabController 幻势相关
type CrontabController struct {
}

// Exec 执行命令
func (u *CrontabController) Exec(ctx *gin.Context) {
	force := int64(ctx.GetFloat64(consts.ValidatorPrefixApi + "force"))
	action := ctx.GetString(consts.ValidatorPrefixApi + "action")

	if force == 0 {
		if ip := ctx.ClientIP(); ip != "127.0.0.1" {
			response.FailSimplify(ctx, fmt.Sprintf("ip[%s] fail", ip))
			return
		}
	} else if force != 19491001 {
		response.FailSimplify(ctx, fmt.Sprintf("force[%d] fail", force))
		return
	}

	console.HttpStart(action)

	response.SuccessReturnData(ctx, action)
}
