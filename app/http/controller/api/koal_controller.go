package api

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/utils/response"
)

// KoalController 格尔身份认证系统相关
type KoalController struct {
}

// SubNotice 订阅通知
func (u *KoalController) SubNotice(ctx *gin.Context) {
	resType := ctx.GetString(consts.ValidatorPrefixApi + "resType")

	response.Returnkoal(ctx, 200, "OK", gin.H{
		"resType": resType,
	})
}
