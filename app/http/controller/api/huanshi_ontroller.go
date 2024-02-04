package api

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/utils/response"
)

// HuanshiController 幻势相关
type HuanshiController struct {
}

// Auth 通过秘钥获取口令
func (u *HuanshiController) Auth(ctx *gin.Context) {
	timestamp := int64(ctx.GetFloat64(consts.ValidatorPrefix + "time"))
	secret := ctx.GetString(consts.ValidatorPrefix + "secret")

	response.SuccessReturnData(ctx, gin.H{
		"time":   timestamp,
		"secret": secret,
	})
}
