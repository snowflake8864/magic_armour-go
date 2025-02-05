package api

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/organization/organization_models"
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
	hsAccessFact := huanshi_platform.CreateCallbackFactory()
	if isOk, errMsg := hsAccessFact.CheckSecret(deviceId, timestamp, secret); !isOk {
		response.FailRoutine(ctx, 10001, errMsg)
		return
	}

	//验证企业
	if model, err := organization_models.CreateOrganizationFactory().GetItemByCode(deviceId); err != nil {
		response.FailRoutine(ctx, 10003, "数据错误")
		return
	} else if model.Id == 0 {
		response.FailRoutine(ctx, 10003, "数据不存在")
		return
	} else {
		variable.ZapLog.Sugar().Info("幻势回调企业编号: ", model.Id)
	}

	//生成口令
	if token, err := hsAccessFact.GenerateToken(deviceId, secret); err != nil {
		response.FailRoutine(ctx, 10002, err.Error())
		return
	} else {
		response.SuccessReturnData(ctx, token)
		return
	}
}
