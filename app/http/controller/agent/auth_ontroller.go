package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/veda_other/huanshi_service"
	"goskeleton/app/utils/response"
)

// AuthController 幻势相关
type AuthController struct {
}

// Auth 通过秘钥获取口令
func (u *AuthController) Auth(ctx *gin.Context) {
	timestamp := int64(ctx.GetFloat64(consts.ValidatorPrefix + "time"))
	secret := ctx.GetString(consts.ValidatorPrefix + "secret")

	//获取备标识码
	dictCacheFact := dictionary_service.CreateCacheFactory()
	deviceId := dictCacheFact.GetValById(enum_dictionary.SystemDeviceId)
	if deviceId == "" {
		response.FailRoutine(ctx, 10000, "设备标识码没有初始化")
	}

	//验证秘钥
	hsAccessFact := huanshi_service.CreateAccessFactory(deviceId)
	if isOk, errMsg := hsAccessFact.CheckSecret(timestamp, secret); !isOk {
		response.FailRoutine(ctx, 10001, errMsg)
	}

	//生成口令
	token, err := hsAccessFact.GenerateToken(secret)
	if err != nil {
		response.FailRoutine(ctx, 10002, err.Error())
	}
	response.SuccessReturnData(ctx, token)
}
