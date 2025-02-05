package api

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/response"
)

// HuanshiController 幻势相关
type HuanshiController struct {
}

// Auth 通过秘钥获取口令
func (u *HuanshiController) Auth(ctx *gin.Context) {
	timestamp := int64(ctx.GetFloat64(consts.ValidatorPrefixApi + "time"))
	secret := ctx.GetString(consts.ValidatorPrefixApi + "secret")
	deviceId := ctx.GetString(consts.ValidatorPrefixApi + "deviceId")

	//验证秘钥
	hsAccessFact := huanshi_platform.CreateCallbackFactory(deviceId)
	if isOk, errMsg := hsAccessFact.CheckSecret(timestamp, secret); !isOk {
		response.FailRoutine(ctx, 10001, errMsg)
		return
	}

	//生成口令
	token, err := hsAccessFact.GenerateToken(secret)
	if err != nil {
		response.FailRoutine(ctx, 10002, err.Error())
		return
	}
	response.SuccessReturnData(ctx, token)
}
